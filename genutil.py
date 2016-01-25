# encoding: utf8

import re
import logging
import clang
import clang.cindex as clidx

FORMAT = '%(asctime)-15s %(filename)s:%(lineno)d %(funcName)s %(message)s'
LOGLEVEL = logging.DEBUG
# LOGLEVEL = logging.ERROR
logging.basicConfig(format=FORMAT, level=LOGLEVEL)
glog = logging.getLogger()


class GenUtil(object):
    basecls = {}
    methods = {}  # cls displayname => class cursor
    signals = {}
    clstokens = {}  # clsname => tokens[]
    inline_methods = {}  # clsname => mangled method name[]
    ticlasses = {}  # template instantiate class

    def __init__(self):
        self.conflib = clang.cindex.conf.lib
        return

    def get_code_file(self, cursor):
        loc = cursor.location
        code_file = loc.file.name.split('/')[-1].split('.')[0]
        return code_file

    # like QtCore
    def get_decl_module(self, cursor):
        loc = cursor.location
        decl_module = loc.file.name.split('/')[-2]
        return decl_module

    # like core
    def get_decl_mod(self, cursor):
        loc = cursor.location
        decl_module = loc.file.name.split('/')[-2][2:].lower()

        # 自测
        if decl_module not in ['core', 'gui', 'widgets', 'network', 'qml', 'quick']:
            raise 'unknown module: %s, %s' % (decl_module, cursor.spelling)

        return decl_module

    # TODO 好像有点bug。
    # QListData会推导出基类是NotIndirectLayout。而实际上QListData没有基类。
    def get_base_class(self, cursor):
        if cursor.displayname in GenUtil.basecls:
            return GenUtil.basecls[cursor.displayname]

        bases = []
        for x in cursor.walk_preorder():
            if x.semantic_parent is not None and \
               x.semantic_parent.kind != clidx.CursorKind.TRANSLATION_UNIT:
                break  # 已经遍历进类内部了，应该需要跳出执行。
            if x.kind == clidx.CursorKind.CXX_BASE_SPECIFIER:
                xdef = x.get_definition()
                if xdef is None:
                    # print(x.kind, x.spelling, cursor.spelling)
                    continue
                decl = x.get_definition().type.get_declaration()

                # if decl.semantic_parent is None:
                if decl.kind == clidx.CursorKind.NO_DECL_FOUND:
                    if xdef.kind != clidx.CursorKind.CLASS_TEMPLATE:
                        print(decl.kind, decl.spelling, x.kind, x.spelling, xdef.kind,
                              cursor.spelling, xdef.location, xdef)
                        raise 'wtf'
                    bases.append(xdef)
                elif decl.semantic_parent.kind == clidx.CursorKind.TRANSLATION_UNIT:
                    bases.append(decl)
                else: break  # 提前跳出结束执行
        GenUtil.basecls[cursor.displayname] = bases
        return bases

    def is_qobject_subclass(self, cursor):
        bases = self.get_base_class(cursor)
        if len(bases) > 0:
            if bases[0].spelling == 'QObject': return True
            else: return self.is_qobject_subclass(bases[0])
        return False

    def get_methods(self, class_cursor):
        if class_cursor.spelling in GenUtil.methods:
            return GenUtil.methods[class_cursor.displayname]

        method_names = {}

        for m in class_cursor.get_children():
            # print(m.kind, m.spelling)
            # TODO va_list type
            # if self.check_skip_method(m): continue
            method_name = m.spelling
            mangled_name = m.mangled_name
            # 这里不需要检测是否是definition，因为这是在类内部的，全部要考虑的
            if m.kind == clidx.CursorKind.CONSTRUCTOR:  # and not m.is_definition():
                method_names[mangled_name] = m
            elif m.kind == clidx.CursorKind.DESTRUCTOR:  # and not m.is_definition():
                method_names[mangled_name] = m
            elif m.kind == clidx.CursorKind.CXX_METHOD:  # and not m.is_definition():
                method_names[mangled_name] = m

        GenUtil.methods[class_cursor.displayname] = method_names
        return method_names

    def get_inst_methods(self, class_cursor, inst_class_cursor):
        tic = self.isTempInstClass(inst_class_cursor)
        template_methods = self.get_methods(class_cursor)

        tip = 'P' if '*' in tic[1] else ''
        caname = tic[1].replace('*', '').strip()
        tisym = 'I%s%d%s' % (tip, len(caname), caname)
        if len(tic[2].strip()) > 0:
            tip = 'P' if '*' in tic[2] else ''
            caname = tic[2].replace('*', '').strip()
            tisym += '%s%d%s' % (tip, len(caname), caname)

        timethods = {}
        for mangled_name in template_methods:
            method_cursor = template_methods[mangled_name]
            ipos = mangled_name.index(inst_class_cursor.spelling)
            ipos = ipos + len(inst_class_cursor.spelling)
            ti_mangled_name = mangled_name[0:ipos] + tisym + mangled_name[ipos:]
            timethods[ti_mangled_name] = template_methods[mangled_name]
            # print(ti_mangled_name, mangled_name, method_cursor.get_num_template_arguments())

        return timethods

    def get_signals(self, cursor):
        if cursor.spelling in GenUtil.signals:
            return GenUtil.signals[cursor.displayname]

        # for it in cursor.walk_preorder():
        #    print(it.kind, it.spelling, it.displayname)
        methods = self.get_methods(cursor)
        signals = {}
        insig = False
        for tk in cursor.get_tokens():
            # print(tk.kind, tk.spelling, tk.cursor.kind)
            if tk.kind == clidx.TokenKind.IDENTIFIER and tk.spelling == 'Q_SIGNALS':
                insig = True
                continue
            if tk.kind == clidx.TokenKind.IDENTIFIER \
               and tk.cursor.kind == clidx.CursorKind.CXX_ACCESS_SPEC_DECL:
                if insig is True:
                    insig = False
                    break
            if tk.kind == clidx.TokenKind.IDENTIFIER and tk.spelling == 'Q_SLOTS':
                if insig is True:
                    insig = False
                    break

            if insig is True and tk.kind == clidx.TokenKind.IDENTIFIER \
               and tk.cursor.kind == clidx.CursorKind.CXX_METHOD:
                # print('got a signal:', tk.spelling, tk.cursor.displayname)
                real_method = methods[tk.cursor.mangled_name]
                signals[tk.cursor.mangled_name] = real_method
                # signals.append(tk.cursor)  # 这种方式拿到的method_cursor有问题

        # print('got signals:', len(signals), signals)
        GenUtil.signals[cursor.displayname] = signals
        return signals

    # qt中inline方法的5种实现方式。
    def get_inline_methods(self, cursor):
        tokens = []
        if cursor.displayname not in GenUtil.clstokens:
            for token in cursor.get_tokens():
                tokens.append(token)
            GenUtil.clstokens[cursor.displayname] = tokens
        else:
            tokens = GenUtil.clstokens[cursor.displayname]

        def care_cond(token):
            if cursor.displayname == 'QModelIndex' and token.cursor.spelling == 'QModelIndex':
                return False
            return False

        inline_methods = []
        if cursor.displayname not in GenUtil.inline_methods:
            all_methods = self.get_methods(cursor)
            pidx = -1
            bidx = 0
            for token in tokens:
                pidx += 1
                if token.cursor.kind == clidx.CursorKind.CONSTRUCTOR \
                   or token.cursor.kind == clidx.CursorKind.CXX_METHOD:
                    if care_cond(token):
                        for tk in tokens[pidx-5:pidx+5]:
                            print(tk.kind, tk.spelling)
                    bidx = pidx
                    while bidx > 0 and tokens[bidx].spelling not in [';', '}', 'Q_DECL_CONSTEXPR']:
                        if tokens[bidx].spelling == 'inline':
                            if care_cond(token):
                                print('found inline 1')
                            inline_methods.append(token.cursor.mangled_name)
                            break
                        bidx -= 1
                    if token.cursor.is_definition():
                        if care_cond(token):
                            print('found inline 2')
                        inline_methods.append(token.cursor.mangled_name)
                    else:
                        if token.cursor.mangled_name not in all_methods:
                            # print('whyyyy,', token.cursor.mangled_name)
                            pass
                        else:
                            mc = all_methods[token.cursor.mangled_name]
                            defn = mc.get_definition()
                            if defn is not None and not mc.is_definition():
                                if care_cond(token):
                                    print('found inline 3', defn.kind, defn.spelling, defn.displayname, defn.location)
                                    for tk in defn.get_tokens():
                                        print(111, tk.kind, tk.spelling)
                                inline_methods.append(token.cursor.mangled_name)
                    # bidx = pidx
                    # while bidx < len(tokens) and tokens[bidx].spelling not in [';']:  # 有方法体,但没inline标识
                    #     if tokens[bidx].spelling == '{':
                    #         inline_methods.append(token.cursor.mangled_name)
                    #         break
                    #     bidx += 1

            GenUtil.inline_methods[cursor.displayname] = inline_methods
        else:
            inline_methods = GenUtil.inline_methods[cursor.displayname]
        return inline_methods

    def get_unique_signals(self, cursor):
        signals = self.get_signals(cursor)
        usignals = {}
        for key in signals:
            sigmth = signals[key]
            usignals[sigmth.spelling] = sigmth

        return usignals

    def is_private_signal(self, method_cursor):
        if 'QPrivateSignal' in method_cursor.displayname: return True
        return False

    # 还要验证基类是否有纯虚方法
    def isAbstractClass(self, cursor):
        for m in cursor.get_children():
            pv = self.conflib.clang_CXXMethod_isPureVirtual(m)
            if pv: return True

        cname = cursor.spelling
        if cname.startswith('QAbstract'): return True

        return False

    def is_pure_virtual_method(self, method_cursor):
        tokens = []
        for token in method_cursor.get_tokens():
            tokens.append(token.spelling)
        if ''.join(tokens[-3:]) == '=0;':
            return True
        return False

    # 是这Qt特有的吗？
    # how
    def isDisableCopy(self, cursor):
        if cursor.spelling == 'QDataStream': return True
        return False

    # deleted是针对C++的构造函数
    # how
    def isDeleted(self, cursor):
        return False

    def isCopyCtor(self, cursor):
        return False

    def isMoveCtor(self, cursor):

        return False

    def hasDefaultCtor(self, cursor):
        return False

    def isqtloc(self, cursor):
        return cursor.location.file.name.startswith('/usr/include/qt')

    def flat_template_name(self, name):
        flat_class_name = name.replace('<', 'L') \
                              .replace('>', 'G') \
                              .replace(':', '_').replace('*', 'P') \
                              .replace(',', '_').replace(' ', 'E')
        return flat_class_name

    def isTempInstClass(self, cursor):
        exp = '^(Q[A-Z].+)\<([^,]*)[, ]*?([^,]+)?[, ]*?([^,]+)?\>'
        res = re.findall(exp, cursor.type.spelling)
        # print(123456, res, cursor.type.spelling, cursor.spelling)
        if len(res) > 0:
            return res[0]
        return None

    def getTempInstClass(self, cursor):
        c = self.conflib.clang_getSpecializedCursorTemplate(cursor)
        return c

    def getNumTempArgs(self, cursor):
        num = 0
        for c in cursor.get_children():
            if c.kind == clidx.CursorKind.TEMPLATE_TYPE_PARAMETER:
                num += 1
        return num

    def getTempArg(self, cursor, idx):
        num = 0
        for c in cursor.get_children():
            if c.kind == clidx.CursorKind.TEMPLATE_TYPE_PARAMETER:
                if num == idx: return c
                num += 1
        return None

    pass


# 可以多点写入的代码编辑类
# 支持多点写入
# 支持前身写入
# 支持唯一写入
class CodePaper:
    def __init__(self):
        self.code = ''
        self.insert_points = {}  # name => [codes]
        self.export_times = 0
        self.newline = '\n'
        return

    def addPoint(self, name):
        if name not in self.insert_points:
            self.insert_points[name] = []
        return

    def hasPoint(self, name):
        if name in self.insert_points: return True
        return False

    def allPoints(self):
        return self.insert_points.keys()

    def append(self, name, code):
        if not self.hasPoint(name): self.addPoint(name)
        self.insert_points[name].append(code)
        return

    def appendUnique(self, name, code):
        if not self.hasPoint(name): self.addPoint(name)
        if code not in self.insert_points[name]:
            self.insert_points[name].append(code)
        return

    def prepend(self, name, code):
        if not self.hasPoint(name): self.addPoint(name)
        self.insert_points[name].insert(0, code)
        return

    def prependUnique(self, name, code):
        if not self.hasPoint(name): self.addPoint(name)
        if code not in self.insert_points[name]:
            self.insert_points[name].insert(0, code)
        return

    def AP(self, name, code): return self.append(name, code)

    def APU(self, name, code): return self.appendUnique(name, code)

    def PP(self, name, code): return self.prepend(name, code)

    def PPU(self, name, code): return self.prependUnique(name, code)

    def getPoint(self, name):
        return self.newline.join(self.insert_points[name])

    def removePoint(self, name):
        codes = self.insert_points.pop(name)
        return self.newline.join(codes)

    def removeLine(self, name, code):
        if code in self.insert_points[name]:
            self.insert_points[name].remove(code)
        return

    # 按照names给出的顺序合并并导出代码。
    def exportCode(self, names):
        self.export_times += 1
        code = ''
        for name in names:
            code += self.newline.join(self.insert_points[name]) + self.newline
        return code

    def totalLength(self):
        tlen = 0
        for name in self.insert_points.keys():
            for line in self.insert_points[name]:
                tlen += len(line)
        return tlen

    def totalLine(self):
        tline = 0
        for name in self.insert_points.keys():
            tline += len(self.insert_points[name])
        return tline

    def reset(self):
        if self.export_times == 0:
            print('Warning, code maybe not export')
        self.insert_points = {}
        self.export_times = 0
        return

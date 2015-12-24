# encoding: utf8
import os
import time

import clang.cindex as clidx

from genbase import GenerateBase, TestBuilder, GenMethodContext, GenClassContext
from genutil import CodePaper, GenUtil
from typeconv import TypeConvForRust


class GenerateForInlineCXX(GenerateBase):
    def __init__(self):
        super(GenerateForInlineCXX, self).__init__()

        self.modrss = {}  # mod => CodePaper
        self.class_blocks = ['header', 'main', 'use', 'ext', 'body']
        return

    def generateHeader(self, module):
        code_file = module
        code = ''
        # code += "#include <QtCore>\n"
        # code += "#include <QtGui>\n"
        # code += "#include <QtWidgets>\n\n"
        code += "#include <%s.h>\n\n" % (code_file)
        code += "extern \"C\" {\n"
        return code

    def generateFooter(self, module):
        code = ''
        code += "} // end extern \"C\" // %s \n" % (module)
        return code

    def generateCMake(self, module, class_decls):
        code = ''

        for elems in class_decls:
            class_name, cs, methods = elems
            code += "  src/%s/%s.cxx\n" % (module[2:].lower(), class_name.lower())

        return code

    def genpass_init_code_paper(self):
        for key in self.gctx.codes:
            CP = self.gctx.codes[key]
            mod = self.gctx.get_decl_mod_by_path(key)
            code_file = self.gctx.get_code_file_by_path(key)
            CP.AP('header', '// auto generated, do not modify.')
            CP.AP('header', '// created: ' + time.ctime())
            CP.AP('header', '// src-file: ' + key)
            CP.AP('header', '// dst-file: /src/%s/%s.cxx' % (mod, code_file))
            CP.AP('header', '//\n')

            for blk in self.class_blocks:
                CP.addPoint(blk)
                CP.append(blk, "// %s block begin =>" % (blk))
        return

    def genpass_code_header(self):
        for key in self.gctx.codes:
            CP = self.gctx.codes[key]
            code_file = self.gctx.get_code_file_by_path(key)
            CP.AP('header', self.generateHeader(code_file))
            # CP.AP('ext', "#[link(name = \"Qt5Core\")]")
            # CP.AP('ext', "#[link(name = \"Qt5Gui\")]")
            # CP.AP('ext', "#[link(name = \"Qt5Widgets\")]\n")
            # CP.AP('ext', "extern {")

        return

    def genpass_code_endian(self):
        for key in self.gctx.codes:
            CP = self.gctx.codes[key]
            CP.append('header', "}; // <= extern \"C\" block end\n")
            for blk in self.class_blocks:
                CP.append(blk, "// <= %s block end\n" % (blk))
                # if blk == 'ext':
                #     CP.append(blk, "} // <= %s block end\n" % (blk))
                # else:
                #     CP.append(blk, "// <= %s block end\n" % (blk))

        return

    def genpass_class_modef(self):

        for key in self.gctx.classes:
            cursor = self.gctx.classes[key]
            if self.check_skip_class(cursor): continue

            class_name = cursor.spelling
            decl_file = self.gctx.get_decl_file(cursor)
            decl_mod = self.gctx.get_decl_mod(cursor)
            code_file = self.gctx.get_code_file(cursor)
            istpl = self.gctx.is_template(cursor)

            if decl_mod not in self.modrss:
                self.modrss[decl_mod] = CodePaper()
                self.modrss[decl_mod].addPoint('main')
                MP = self.modrss[decl_mod]

            MP = self.modrss[decl_mod]
            # MP.APU('main', "pub mod %s;" % (code_file))
            # MP.APU('main', "pub use self::%s::%s;\n" % (code_file, class_name))
            MP.APU('main', "set(qt5_inline_%s_srcs ${qt5_inline_%s_srcs} src/%s/%s.cxx)" %
                   (decl_mod, decl_mod, decl_mod, code_file))
        return

    def genpass(self):
        self.genpass_init_code_paper()
        self.genpass_code_header()

        print('gen classes...')
        self.genpass_classes()

        print('gen code endian...')
        self.genpass_code_endian()

        print('gen class mod define...')
        self.genpass_class_modef()

        print('gen files...')
        self.genpass_write_codes()
        return

    def genpass_classes(self):
        for key in self.gctx.classes:
            cursor = self.gctx.classes[key]
            if self.check_skip_class(cursor): continue

            class_name = cursor.displayname
            methods = self.gutil.get_methods(cursor)
            bases = self.gutil.get_base_class(cursor)
            base_class = bases[0] if len(bases) > 0 else None
            # self.generateInheritEmulate(cursor, base_class)
            self.generateClass(class_name, cursor, methods, base_class)
            # break
        return

    # def generateClasses(self, module, class_decls):
    #     code = ''

    #     for elems in class_decls:
    #         class_name, cs, methods = elems
    #         tcode = self.generateClass(class_name, cs, methods)
    #         tcode = self.generateHeader(module) + tcode + self.generateFooter(module)
    #         self.write_code(module, class_name.lower(), tcode)
    #         code += tcode

    #     return code

    def generateClass(self, class_name, class_cursor, methods, base_class):

        # CP = self.gctx.getCodePager(class_cursor)

        # ctysz = class_cursor.type.get_size()
        # CP.AP('body', "// class sizeof(%s)=%s\n" % (class_name, ctysz))

        # generate struct of class
        # CP.AP('body', "pub struct %s {\n" % (class_name))
        # CP.AP('body', "  pub qclsinst: *mut c_void,\n")
        # CP.AP('body', "}\n\n")

        # 重载的方法，只生成一次trait
        unique_methods = {}
        for mangled_name in methods:
            cursor = methods[mangled_name]
            isstatic = cursor.is_static_method()
            static_suffix = '_s' if isstatic else ''
            umethod_name = cursor.spelling + static_suffix
            unique_methods[umethod_name] = True

        # dupremove = self.dedup_return_const_diff_method(methods)
        dupremove = []
        # print(444, 'dupremove len:', len(dupremove), dupremove)
        for mangled_name in methods:
            cursor = methods[mangled_name]
            method_name = cursor.spelling
            if self.check_skip_method(cursor):
                # if method_name == 'QAction':
                    #print(433, 'whyyyyyyyyyyyyyy') # no
                    # exit(0)
                continue
            if mangled_name in dupremove:
                # print(333, 'skip method:', mangled_name)
                continue

            ctx = self.createGenMethodContext(cursor, class_cursor, base_class, unique_methods)
            self.generateMethod(ctx)

        return

    def createGenMethodContext(self, method_cursor, class_cursor, base_class, unique_methods):
        ctx = GenMethodContext(method_cursor, class_cursor)
        ctx.unique_methods = unique_methods
        ctx.CP = self.gctx.getCodePager(class_cursor)

        # if ctx.ctor: ctx.method_name_rewrite = 'New%s' % (ctx.method_name)
        # if ctx.dtor: ctx.method_name_rewrite = 'Free%s' % (ctx.method_name[1:])
        if ctx.ctor: ctx.method_name_rewrite = 'New'
        if ctx.dtor: ctx.method_name_rewrite = 'Free'
        if self.is_conflict_method_name(ctx.method_name):
            ctx.method_name_rewrite = ctx.method_name + '_'
        if ctx.static:
            ctx.method_name_rewrite = ctx.method_name + ctx.static_suffix

        class_name = ctx.class_name
        method_name = ctx.method_name

        # ctx.ret_type_name_rs = self.tyconv.Type2RustRet(ctx.ret_type, method_cursor)
        # ctx.ret_type_name_ext = self.tyconv.TypeCXX2RustExtern(ctx.ret_type, method_cursor)

        raw_params_array = self.generateParamsRaw(class_name, method_name, method_cursor)
        raw_params = ', '.join(raw_params_array)

        # trait_params_array = self.generateParamsForTrait(class_name, method_name, method_cursor, ctx)
        # trait_params = ', '.join(trait_params_array)

        call_params_array = self.generateParamsForCall(class_name, method_name, method_cursor)
        call_params = ', '.join(call_params_array)
        if not ctx.static and not ctx.ctor: call_params = ('rsthis.qclsinst, ' + call_params).strip(' ,')

        # extargs_array = self.generateParamsForExtern(class_name, method_name, method_cursor, ctx)
        # extargs = ', '.join(extargs_array)
        # if not ctx.static: extargs = ('qthis: *mut c_void, ' + extargs).strip(' ,')

        ctx.params_cpp = raw_params
        # ctx.params_rs = trait_params
        ctx.params_call = call_params
        # ctx.params_ext = extargs

        # ctx.trait_proto = '%s::%s(%s)' % (class_name, method_name, trait_params)
        ctx.fn_proto_cpp = "  // proto: %s %s %s::%s(%s);" % \
                           (ctx.static_str, ctx.ret_type_name_cpp, ctx.class_name, ctx.method_name, ctx.params_cpp)
        ctx.has_return = self.methodHasReturn(ctx)

        # base class
        ctx.base_class = base_class
        ctx.base_class_name = base_class.spelling if base_class is not None else None
        ctx.has_base = True if base_class is not None else False
        ctx.has_base = base_class is not None

        # aux
        ctx.tymap = TypeConvForRust.tymap

        return ctx

    def createMiniContext(self, cursor, base_class):
        ctx = GenClassContext(cursor)
        ctx.CP = self.gctx.getCodePager(cursor)

        # base class
        ctx.base_class = base_class
        ctx.base_class_name = base_class.spelling if base_class is not None else None
        ctx.has_base = True if base_class is not None else False
        ctx.has_base = base_class is not None

        # aux
        ctx.tymap = TypeConvForRust.tymap
        return ctx

    # def generateClass(self, class_name, cs, methods, base_class):
    #     code = ''

    #     ctysz = cs.type.get_size()
    #     code += "// class sizeof(%s)=%s\n" % (class_name, ctysz)

    #     for mth in methods:
    #         cursor = methods[mth]
    #         if self.check_skip_method(cursor): continue

    #         code += self.generateMethod(class_name, mth, cursor)
    #         # print(111, code)

    #     return code

    def generateMethod(self, ctx):
        class_name = ctx.class_name
        method_name = ctx.method_name
        method_cursor = ctx.cursor
        cursor = method_cursor
        code = ''

        return_type = cursor.result_type
        return_real_type = self.real_type_name(return_type)
        if '::' in return_real_type: return code
        # if self.check_skip_params(cursor): return code

        inner_return = ''
        if cursor.kind == clidx.CursorKind.CONSTRUCTOR or \
           cursor.kind == clidx.CursorKind.DESTRUCTOR:
            # code += " %s(" % (fixmthname)
            code += "void "
            pass
        else:
            return_type_name = self.resolve_swig_type_name(class_name, return_type)
            return_type_name2 = self.hotfix_typename_ifenum_asint(class_name, method_cursor, return_type)
            return_type_name = return_type_name2 if return_type_name2 is not None else return_type_name
            # code += " %s %s(" % (return_type_name, fixmthname)
            code += "%s " % (return_type_name)
            inner_return = 'return' if return_type_name != 'void' else inner_return

        params = self.generateParams(class_name, method_name, method_cursor);
        params = ', '.join(params)
        params = 'void *that, ' + params
        params = params.strip(', ')

        ret_type_name = ctx.ret_type_name_cpp.replace('&', '*') if ctx.ret_type_ref else ctx.ret_type_name_cpp
        ctx.CP.AP('header', ctx.fn_proto_cpp)
        mangled_name = method_cursor.mangled_name
        if return_type.kind == clidx.TypeKind.RECORD:
            code += "%s(%s)\n" % (mangled_name, params)
            ctx.CP.AP('header', "%s* %s(%s)\n" % (ret_type_name, mangled_name, params))
        else:
            code += "%s(%s)\n" % (mangled_name, params)
            ctx.CP.AP('header', "%s %s(%s)\n" % (ret_type_name, mangled_name, params))

        code += "{\n"
        # code += "  fprintf(stderr, \"Do't call this function.\\n\"); exit(-1);\n"
        ctx.CP.AP('header', "{")

        call_params = self.generateParamsForCall(class_name, method_name, method_cursor)
        call_params = ', '.join(call_params)
        code += "  %s *cthat = (%s *)that;\n" % (class_name, class_name)
        ctx.CP.AP('header', "  %s *cthat = (%s *)that;" % (class_name, class_name))
        if cursor.kind == clidx.CursorKind.CONSTRUCTOR:
            code += "  auto _o = new(that) %s(%s);\n" % (method_name, call_params)
            ctx.CP.AP('header', "  auto _o = new(that) %s(%s);" % (method_name, call_params))
        else:
            code += "  %s cthat->%s(%s);\n" % (inner_return, method_name, call_params)
            if ctx.ret_type_ref:
                ctx.CP.AP('header', "  %s &cthat->%s(%s);" % (inner_return, method_name, call_params))
            else:
                if return_type.kind == clidx.TypeKind.RECORD:
                    ctx.CP.AP('header', "  auto recret = cthat->%s(%s);" % (method_name, call_params))
                    ctx.CP.AP('header', "  %s new %s(recret);" % (inner_return, return_type.spelling))
                    # ctx.CP.AP('header', "  %s std::move(cthat->%s(%s));" % (inner_return, method_name, call_params))
                else:
                    ctx.CP.AP('header', "  %s cthat->%s(%s);" % (inner_return, method_name, call_params))

        code += "}\n"
        ctx.CP.AP('header', '}\n')
        return code

    # @return []
    def generateParamsRaw(self, class_name, method_name, method_cursor):
        argv = []
        for arg in method_cursor.get_arguments():
            argelem = "%s %s" % (arg.type.spelling, arg.displayname)
            argv.append(argelem)
        return argv

    # @return []
    def generateParams(self, class_name, method_name, method_cursor):
        idx = 0
        argv = []

        for arg in method_cursor.get_arguments():
            idx += 1
            # print('%s, %s, ty:%s, kindty:%s' % (method_name, arg.displayname, arg.type.spelling, arg.kind))
            # print('arg type kind: %s, %s' % (arg.type.kind, arg.type.get_declaration()))
            # param_line2 = self.restore_param_by_token(arg)
            # print(param_line2)

            type_name = self.resolve_swig_type_name(class_name, arg.type)
            type_name2 = self.hotfix_typename_ifenum_asint(class_name, arg, arg.type)
            type_name = type_name2 if type_name2 is not None else type_name

            arg_name = 'arg%s' % idx if arg.displayname == '' else arg.displayname
            # try fix void (*)(void *) 函数指针
            # 实际上swig不需要给定名字，只需要类型即可。
            if arg.type.kind == clidx.TypeKind.POINTER and "(*)" in type_name:
                argelem = "%s" % (type_name.replace('(*)', '(*%s)' % arg_name))
            else:
                argelem = "%s %s" % (type_name, arg_name)
            argv.append(argelem)

        return argv

    # @return []
    def generateParamsForCall(self, class_name, method_name, method_cursor):
        idx = 0
        argv = []

        for arg in method_cursor.get_arguments():
            idx += 1
            # print('%s, %s, ty:%s, kindty:%s' % (method_name, arg.displayname, arg.type.spelling, arg.kind))
            # print('arg type kind: %s, %s' % (arg.type.kind, arg.type.get_declaration()))
            # param_line2 = self.restore_param_by_token(arg)
            # print(param_line2)

            type_name = self.resolve_swig_type_name(class_name, arg.type)
            type_name2 = self.hotfix_typename_ifenum_asint(class_name, arg, arg.type)
            type_name = type_name2 if type_name2 is not None else type_name

            arg_name = 'arg%s' % idx if arg.displayname == '' else arg.displayname
            argelem = "%s" % (arg_name)
            argv.append(argelem)

        return argv

    def methodHasReturn(self, ctx):
        method_cursor = cursor = ctx.cursor
        class_name = ctx.class_name

        return_type = cursor.result_type

        return_type_name = return_type.spelling
        if ctx.ctor or ctx.dtor: pass
        else:
            return_type_name = self.resolve_swig_type_name(class_name, return_type)
            return_type_name2 = self.hotfix_typename_ifenum_asint(class_name, method_cursor, return_type)
            return_type_name = return_type_name2 if return_type_name2 is not None else return_type_name

        has_return = True
        if return_type_name == 'void': has_return = False
        # if cursor.spelling == 'buttons':
        #     print(666, has_return, return_type_name, cursor.spelling, return_type.kind, cursor.semantic_parent.spelling)
        #     exit(0)
        if '<' in return_type_name: has_return = False
        if "QStringList" in return_type_name: has_return = False
        if "QObjectList" in return_type_name: has_return = False
        if '::' in return_type_name: has_return = False
        if 'QAbstract' in return_type_name: has_return = False
        if 'QMetaObject' in return_type_name: has_return = False
        if 'QOpenGL' in return_type_name: has_return = False
        if 'QGraphics' in return_type_name: has_return = False
        if 'QPlatform' in return_type_name: has_return = False
        if 'QFunctionPointer' in return_type_name: has_return = False
        if 'QTextEngine' in return_type_name: has_return = False
        if 'QTextDocumentPrivate' in return_type_name: has_return = False
        if 'QJson' in return_type_name: has_return = False
        if 'QStringRef' in return_type_name: has_return = False

        if 'internalPointer' in method_cursor.spelling: has_return = False
        if 'rwidth' in method_cursor.spelling: has_return = False
        if 'rheight' in method_cursor.spelling: has_return = False
        if 'utf16' == method_cursor.spelling: has_return = False
        if 'x' == method_cursor.spelling: has_return = False
        if 'rx' == method_cursor.spelling: has_return = False
        if 'y' == method_cursor.spelling: has_return = False
        if 'ry' == method_cursor.spelling: has_return = False
        if class_name == 'QGenericArgument' and method_cursor.spelling == 'data': has_return = False
        if class_name == 'QSharedMemory' and method_cursor.spelling == 'constData': has_return = False
        if class_name == 'QSharedMemory' and method_cursor.spelling == 'data': has_return = False
        if class_name == 'QVariant' and method_cursor.spelling == 'constData': has_return = False
        if class_name == 'QVariant' and method_cursor.spelling == 'data': has_return = False
        if class_name == 'QThreadStorageData' and method_cursor.spelling == 'set': has_return = False
        if class_name == 'QThreadStorageData' and method_cursor.spelling == 'get': has_return = False
        if class_name == 'QChar' and method_cursor.spelling == 'unicode': has_return = False

        return has_return

    def is_conflict_method_name(self, method_name):
        return False
        if method_name in ['match', 'type', 'move']:  # , 'select']:
            return True
        return False

    # @return True | False
    def check_skip_method(self, cursor):
        method_name = cursor.spelling
        if method_name.startswith('operator'):
            # print("Omited operator method: " + mth)
            return True

        if not self.method_is_inline(cursor): return True

        # print('pub:' + str(cursor.access_specifier))
        if cursor.access_specifier == clidx.AccessSpecifier.PUBLIC:
            pass
        if cursor.access_specifier == clidx.AccessSpecifier.PROTECTED:
            return True
            # if cursor.kind == clidx.CursorKind.CONSTRUCTOR or \
            #   cursor.kind == clidx.CursorKind.DESTRUCTOR:
            #    pass
            # else: return True
        if cursor.access_specifier == clidx.AccessSpecifier.PRIVATE:
            return True
            # if cursor.kind == clidx.CursorKind.CONSTRUCTOR or \
            #   cursor.kind == clidx.CursorKind.DESTRUCTOR:
            #    pass
            # else: return True

        istatic = cursor.is_static_method()
        # if istatic is True: return True

        # fix method
        fixmths = ['tr', 'trUtf8', 'qt_metacall', 'qt_metacast', 'data_ptr',
                   'sprintf', 'vsprintf', 'vasprintf', 'asprintf',
                   'entryInfoListcc',]
        if method_name in fixmths: return True
        fixmths_prefix = ['qt_check_for_']
        for p in fixmths_prefix:
            if method_name.startswith(p): return True

        # 实现不知道怎么fix了，已经fix，原来是给clidx.parse中的-I不全，导致找不到类型。
        # fixmths3 = ['setQueryItems']
        # if method_name in fixmths3: return True

        return False

    def check_skip_class(self, class_cursor):
        cursor = class_cursor
        name = cursor.spelling
        dname = cursor.displayname

        if name in ['QTypeInfo']: return True

        # for template
        if self.gctx.is_template(cursor): return True

        def has_template_brother(cursor):
            for key in self.gctx.classes:
                tc = self.gctx.classes[key]
                if tc != cursor and tc.spelling == cursor.spelling and tc.kind == clidx.CursorKind.CLASS_TEMPLATE:
                    return True
            return False

        hastb = has_template_brother(cursor)
        if hastb: return True

        # like QIntegerForSize<1/2/3>
        if '<' in dname: return True

        # if 'QFuture<' in dname:
        #     for it in cursor.walk_preorder():
        #         print(it.kind, it.displayname, it.location)
        #     print(cursor.get_num_template_arguments())
        #     exit(0)

        return False

    def method_is_inline(self, method_cursor):
        for token in method_cursor.get_tokens():
            if token.spelling == 'inline':
                parent = method_cursor.semantic_parent
                # print(111, method_cursor.spelling, parent.spelling)
                return True
        return False

    # def hotfix_typename_ifenum_asint(self, class_name, arg):
    def hotfix_typename_ifenum_asint(self, class_name, token_cursor, atype):
        type_name = self.resolve_swig_type_name(class_name, atype)
        # if type_name not in ('int', 'int *', 'const int &'): return None
        type_name_segs = type_name.split(' ') 
        if 'int' not in type_name_segs: return None

        tokens = []
        for token in token_cursor.get_tokens():
            tokens.append(token.spelling)
            tkcursor = token.cursor

        # 为什么tokens是空呢，是不能识别的？
        if len(tokens) == 0: return None
        # TODO 全部使用replace方式，而不是这种每个符号的处理
        while tokens[0] in ['const', 'inline']:
            tokens = tokens[1:]

        firstch = tokens[0][0:1]
        if firstch.upper() == firstch and firstch != 'Q':
            print('Warning fix enum-as-int:', type_name, '=> %s::' % class_name, tokens[0])
            return '%s::%s' % (class_name, tokens[0])

        if len(tokens) < 3: return None
        if firstch.upper() == firstch and firstch == 'Q' and tokens[1] == '::':
            print('Warning fix enum-as-int2:', type_name, '=> %s::' % class_name, tokens[2])
            return '%s::%s' % (tokens[0], tokens[2])

        # like QtMsgType
        if firstch.upper() == firstch and firstch == 'Q' and tokens[0][0:2] == 'Qt':
            print('Warning fix enum-as-int3:', type_name, '=> ', tokens[0])
            return '%s' % (tokens[0])

        if firstch.upper() == firstch and firstch == 'Q' and tokens[0][1:2].lower() == tokens[0][1:2]:
            print('Warning fix enum-as-int4:', type_name, '=> ', type_name.replace('int', tokens[0]))
            return '%s' % (type_name.replace('int', tokens[0]))

        # like qint64...
        if firstch.lower() == firstch and tokens[0][0:1] == 'q' and '*' in type_name:
            print('Warning fix qint*-as-int5:', type_name, '=> ', tokens[0])
            return '%s %s' % (tokens[0], tokens[1])

        return None

    def real_type_name(self, atype):
        type_name = atype.spelling

        if atype.kind == clidx.TypeKind.TYPEDEF:
            # print('underlying type: %s' % atype.get_declaration().underlying_typedef_type.spelling)
            # print('underlying type: %s' % arg.type.underlying_typedef_type.spelling)
            type_name = atype.get_declaration().underlying_typedef_type.spelling
            if type_name.startswith('QFlags<'):
                type_name = type_name[7:len(type_name) - 1]

        return type_name

    # @return str
    def resolve_swig_type_name(self, class_name, atype):
        type_name = atype.spelling
        if type_name in ['QFunctionPointer', 'CategoryFilter',
                         'EasingFunction']:
            type_bclass = atype.get_declaration().semantic_parent
            # if type_name.startswith('Q'):
            # 全局定义的，不需要前缀
            if type_bclass.kind == clidx.CursorKind.TRANSLATION_UNIT: pass
            else: type_name = '%s::%s' % (type_bclass.spelling, type_name)
        else:
            type_name = self.real_type_name(atype)

            # QTextStreamManipulator(void (QTextStream::*)(int) m, int a);
            # int registerNormalizedType(const ::QByteArray & normalizedTypeName, void * destructor, void *(*)(void *, const void *) constructor, int size, QMetaType::TypeFlags flags, const QMetaObject * metaObject);
            # qreal (*)(qreal) customType();
            # if type_name == 'void (*)(void *)':
            #    type_name = "void *"

        return type_name

    def write_cmake_code(self, module, fname, code):
        fpath = "CMakeLists.txt"
        f = os.open(fpath, os.O_CREAT | os.O_TRUNC | os.O_RDWR)
        os.write(f, code)
        os.close(f)
        return

    def genpass_write_codes(self):
        for key in self.gctx.codes:
            cp = self.gctx.codes[key]
            code = cp.exportCode(self.class_blocks)

            mod = self.gctx.get_decl_mod_by_path(key)
            fname = self.gctx.get_code_file_by_path(key)
            if mod not in ['core', 'gui', 'widgets', 'network', 'dbus']:
                print('Omit unknown mod code...:', mod, fname, key)
                continue

            self.write_code(mod, fname, code)
            # self.write_file(fpath, code)

        # class mod define
        # self.write_modrs(module, self.MP.exportCode(['main']))
        for mod in self.modrss:
            cp = self.modrss[mod]
            code = cp.exportCode(['main'])
            lines = cp.totalLine()
            print('write mod.cmake:', mod, len(code), lines)
            self.write_modrs(mod, code)
        return

    def write_code(self, mod, fname, code):
        # mod = 'core'
        # fpath = "src/core/%s.rs" % (fname)
        fpath = "src/%s/%s.cxx" % (mod, fname)
        self.write_file(fpath, code)
        return

    # TODO dir is exists
    def write_file(self, fpath, code):
        f = os.open(fpath, os.O_CREAT | os.O_TRUNC | os.O_RDWR)
        os.write(f, code)
        os.close(f)

        return

    def write_modrs(self, mod, code):
        fpath = "src/%s/mod.cmake" % (mod)
        self.write_file(fpath, code)
        return

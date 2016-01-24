# encoding: utf8
import os
import time

import clang.cindex as clidx

from genbase import GenerateBase, TestBuilder, GenMethodContext, GenClassContext
from genutil import CodePaper, GenUtil
from typeconv import TypeConvForRust
from genfilter import GenFilterInc


class GenerateForInc(GenerateBase):
    def __init__(self):
        super(GenerateForInc, self).__init__()

        self.gfilter = GenFilterInc()
        self.modrss = {}  # mod => CodePaper
        self.class_blocks = ['header', 'main', 'use', 'ext', 'body']
        return

    def generateHeader(self, module):
        code_file = module
        code = ''
        # code += "#include <QtCore>\n"
        # code += "#include <QtGui>\n"
        # code += "#include <QtWidgets>\n\n"
        code += "#include <qatomic.h>\n"  # for fix qatomic_x86.cxx's compile
        code += "#include <qstring.h>\n"  # for fix qbytearray.cxx's compile
        code += "#include <qfuture.h>\n"  # for fix qfutureinterface.cxx's compile
        code += "#include <qpoint.h>\n"   # for fix qeasingcurve.cxx's compile
        code += "#include <qurl.h>\n"   # for fix qmetadata.cxx's compile
        code += "#include <qopengl.h>\n"   # for fix qmetadata.cxx's compile

        code += "#include <%s.h>\n\n" % (code_file)
        # code += "extern \"C\" {\n"
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
            CP.AP('main', 'void __keep_%s_inline_symbols() {' % (code_file))

        return

    def genpass_code_endian(self):
        for key in self.gctx.codes:
            CP = self.gctx.codes[key]
            # CP.append('header', "}; // <= extern \"C\" block end\n")
            CP.append('main', "} // <= main block end\n")
            for blk in self.class_blocks:
                CP.append(blk, "// <= %s block end\n" % (blk))

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
            MP.APU('main', "set(qt5_inline_%s_srcs ${qt5_inline_%s_srcs} src/%s/%s.cxx)" %
                   (decl_mod, decl_mod, decl_mod, code_file))
        return

    def genpass(self):
        self.genpass_init_code_paper()
        self.genpass_code_header()

        print('gen classes...')
        self.genpass_classes()

        print('gen instantiate classes...' + str(len(self.gutil.ticlasses)))

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
            self.generateClass(class_name, cursor, methods, base_class)
            # break

        return

    def generateClass(self, class_name, class_cursor, methods, base_class):
        # 重载的方法，只生成一次trait
        unique_methods = {}
        for mangled_name in methods:
            cursor = methods[mangled_name]
            isstatic = cursor.is_static_method()
            static_suffix = '_s' if isstatic else ''
            umethod_name = cursor.spelling + static_suffix
            unique_methods[umethod_name] = True

        # 生成计算类大小的方法
        mctx = self.createMiniContext(class_cursor, base_class)
        self.generateClassSize(mctx)
        self.generateSlotProxy(mctx)

        for mangled_name in methods:
            cursor = methods[mangled_name]
            if self.check_skip_method(cursor):
                continue
            # if not self.method_is_inline(cursor):
            #    continue
            if self.check_skip_params(cursor):
                continue

            ctx = self.createGenMethodContext(cursor, class_cursor, base_class, unique_methods)
            if cursor.kind == clidx.CursorKind.CONSTRUCTOR \
               or cursor.kind == clidx.CursorKind.DESTRUCTOR:
                self.generateCtors(ctx)
            else:
                self.generateMethod(ctx)
                if ctx.isinline:
                    # 生成inline部分symbole
                    self.generateMethodInline(ctx)

        return

    def createGenMethodContext(self, method_cursor, class_cursor, base_class, unique_methods):
        ctx = GenMethodContext(method_cursor, class_cursor)
        ctx.unique_methods = unique_methods
        ctx.CP = self.gctx.getCodePager(class_cursor)

        if ctx.ctor: ctx.method_name_rewrite = 'New'
        if ctx.dtor: ctx.method_name_rewrite = 'Free'
        if self.is_conflict_method_name(ctx.method_name):
            ctx.method_name_rewrite = ctx.method_name + '_'
        if ctx.static:
            ctx.method_name_rewrite = ctx.method_name + ctx.static_suffix

        ctx.isinline = self.method_is_inline(method_cursor)
        ctx.isabstract = self.gutil.isAbstractClass(class_cursor)

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
                           (ctx.static_str, ctx.ret_type_name_cpp, ctx.full_class_name, ctx.method_name, ctx.params_cpp)
        ctx.has_return = self.methodHasReturn(ctx)
        ctx.need_return = self.methodNeedReturn(ctx)

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

    def generateSlotProxy(self, mctx):
        ctx = mctx

        isqobject = self.gutil.is_qobject_subclass(ctx.cursor)
        if isqobject is False: return

        # base_slot_proxy_class_name = ctx.base_class_name
        signals = self.gutil.get_signals(ctx.cursor)

        def gen_proto_line(mth, fptr=False):
            argv = ['void* rsfptr'] if fptr else []
            idx = 0
            for arg in mth.get_arguments():
                full_tyname = arg.type.spelling
                tydecl = arg.type.get_declaration()
                if tydecl is not None and tydecl.semantic_parent is not None:
                    # print(tydecl.semantic_parent.spelling, tydecl.semantic_parent.kind)
                    if tydecl.semantic_parent.spelling == mth.semantic_parent.spelling \
                       and '::' not in full_tyname:
                        full_tyname = '%s::%s' % (mth.semantic_parent.spelling, full_tyname)
                        # print(arg.type.spelling, '==>', full_tyname)
                argv.append('%s arg%s' % (full_tyname, idx))
                idx += 1
            return ', ' .join(argv)

        def gen_call_line(mth):
            argv = ['this->rsfptr']
            for arg in mth.get_arguments():
                argv.append('arg%s' % (len(argv)-1))
            return ', ' .join(argv)

        ctx.CP.AP('body', '// %s_SlotProxy here' % (ctx.class_name))
        ctx.CP.AP('body', 'class %s_SlotProxy : public QObject' % (ctx.flat_class_name))
        ctx.CP.AP('body', '{')
        ctx.CP.AP('body', '  Q_OBJECT;')
        ctx.CP.AP('body', 'public:')
        ctx.CP.AP('body', '   %s_SlotProxy():QObject(){}' % (ctx.flat_class_name))
        ctx.CP.AP('body', '')

        ### signals
        for key in signals:
            sigmth = signals[key]
            if '<' in sigmth.displayname: continue
            if self.gutil.is_private_signal(sigmth): continue
            proto_line = gen_proto_line(sigmth)
            proto_line_fptr = gen_proto_line(sigmth, True)
            ctx.CP.AP('body', 'public slots:')
            ctx.CP.AP('body', '  // %s' % (sigmth.displayname))
            ctx.CP.AP('body', '  void slot_proxy_func_%s(%s);' % (sigmth.mangled_name, proto_line))
            ctx.CP.AP('body', 'public:')
            ctx.CP.AP('body', '  void (*slot_func_%s)(%s) = NULL;' % (sigmth.mangled_name, proto_line_fptr))

        ctx.CP.AP('body', 'public: void* rsfptr = NULL;')
        ctx.CP.AP('body', '};')

        code_mod = self.gctx.get_decl_mod(ctx.cursor)
        code_file = self.gctx.get_code_file(ctx.cursor)
        ctx.CP.removeLine('body', '#include \"src/%s/%s.moc\"' % (code_mod, code_file))
        ctx.CP.AP('body', '#include \"src/%s/%s.moc\"' % (code_mod, code_file))
        ctx.CP.AP('body', '')

        ctx.CP.AP('body', 'extern \"C\" {')
        ctx.CP.AP('body', '  %s_SlotProxy* %s_SlotProxy_new()' % (ctx.flat_class_name, ctx.flat_class_name))
        ctx.CP.AP('body', '  {')
        ctx.CP.AP('body', '    return new %s_SlotProxy();' % (ctx.flat_class_name))
        ctx.CP.AP('body', '  }')
        ctx.CP.AP('body', '};')
        ctx.CP.AP('body', '')

        for key in signals:
            sigmth = signals[key]
            if '<' in sigmth.displayname: continue
            if self.gutil.is_private_signal(sigmth): continue
            proto_line = gen_proto_line(sigmth)
            call_line = gen_call_line(sigmth)
            ctx.CP.AP('body', 'void %s_SlotProxy::slot_proxy_func_%s(%s) {'
                      % (ctx.class_name, sigmth.mangled_name, proto_line))
            ctx.CP.AP('body', '  if (this->slot_func_%s != NULL) {' % (sigmth.mangled_name))
            ctx.CP.AP('body', '    // do smth...')
            ctx.CP.AP('body', '    this->slot_func_%s(%s);' % (sigmth.mangled_name, call_line))
            ctx.CP.AP('body', '  }')
            ctx.CP.AP('body', '}')

            ctx.CP.AP('body', 'extern \"C\"')
            ctx.CP.AP('body', 'void* %s_SlotProxy_connect_%s(QObject* sender, void* ffifptr, void* rsfptr){'
                      % (ctx.flat_class_name, sigmth.mangled_name))
            ctx.CP.AP('body', '  auto that = new %s_SlotProxy();' % (ctx.flat_class_name))
            ctx.CP.AP('body', '  that->rsfptr = rsfptr;')
            ctx.CP.AP('body', '  that->slot_func_%s = (decltype(that->slot_func_%s))ffifptr;'
                      % (sigmth.mangled_name, sigmth.mangled_name))
            # 无法使用C++11的connect方式，有可能重载的方法，不适用。
            ctx.CP.AP('body', '  QObject::connect((%s*)sender, SIGNAL(%s), that, SLOT(slot_proxy_func_%s(%s)));'
                      % (ctx.class_name, sigmth.displayname, sigmth.mangled_name, proto_line))
            ctx.CP.AP('body', '  return that;')
            ctx.CP.AP('body', '}')

            ctx.CP.AP('body', 'extern \"C\"')
            ctx.CP.AP('body', 'void %s_SlotProxy_disconnect_%s(%s_SlotProxy* that) {'
                      % (ctx.flat_class_name, sigmth.mangled_name, ctx.flat_class_name))
            ctx.CP.AP('body', '  that->disconnect();')
            ctx.CP.AP('body', '  delete that;')
            ctx.CP.AP('body', '}')
            ctx.CP.AP('body', '')

        return

    def generateClassSize(self, mctx):
        # 获取类大小的封装，clang.py获取的类大小不对，如果有clang.cpp应该能够获取到正确值吧

        ctx = mctx

        # 类内类处理
        full_class_name = ctx.full_class_name
        symbol_name = full_class_name.replace('<', '_').replace('>', '_').replace(' ', '_') \
                      .replace('const', '').replace('*', '_').replace(':', '_').replace(',', '_')

        ctx.CP.AP('use', 'extern "C"')
        ctx.CP.AP('use', 'int %s_Class_Size()' % (symbol_name))
        ctx.CP.AP('use', '{')
        ctx.CP.AP('use', '  return sizeof(%s);' % (full_class_name))
        ctx.CP.AP('use', '}\n')

        return

    # 生成所有的构造方法
    def generateCtors(self, ctx):
        method_cursor = ctx.cursor
        method_name = ctx.method_name
        mangled_name = ctx.mangled_name

        # 不能实例化
        # TODO 使用ispure做准确判断
        # 移动到上面

        if method_cursor.kind == clidx.CursorKind.CONSTRUCTOR:
            if ctx.isinline:
                self.generateCtorAllocInline(ctx)
            else:
                self.generateCtorAlloc(ctx)

        if method_cursor.kind == clidx.CursorKind.DESTRUCTOR:
            if ctx.isinline:
                self.generateDtorDeleteInline(ctx)
            else:
                self.generateDtorDelete(ctx)

        return

    # 重新生成新的ctor封装，在这里计算类大小，并分配空间，返回生成的对象
    def generateCtorAlloc(self, ctx):

        class_name = ctx.class_name
        method_name = ctx.method_name
        method_cursor = ctx.cursor

        params_arr = self.generateParams(class_name, method_name, method_cursor)
        params = ', '.join(params_arr)

        # 类内类处理
        full_class_name = ctx.full_class_name

        idx = 0
        argv = []
        prmv = []
        for arg in ctx.cursor.get_arguments():
            idx += 1
            # self.generateParamDeclExpr(ctx, arg, idx)
            # argv.append('arg%s' % (idx))
            prmc = self.generateParamForCall(arg, idx)
            argv.append(prmc)
            prme = self.generateParamForDecl(arg, idx)
            prmv.append('%s' % (prme))
        args = ',\n'.join(argv)
        prms = ',\n'.join(prmv)

        ctx.CP.AP('ext', '// %s' % (str(ctx.cursor.location)))
        ctx.CP.AP('ext', '// %s' % (ctx.fn_proto_cpp))
        ctx.CP.AP('ext', 'extern "C"')
        self.generateReturnDecl(ctx)
        ctx.CP.AP('ext', 'C%s(%s) {' % (ctx.mangled_name, prms))
        if ctx.isabstract:
            ctx.CP.AP('ext', '  // auto ret = new %s(%s);' % (ctx.full_class_name, args))
        else:
            ctx.CP.AP('ext', '  auto ret = new %s(%s);' % (ctx.full_class_name, args))
            self.generateReturnImpl(ctx)
        ctx.CP.AP('ext', '}')

        return

    # 重新生成新的ctor封装，在这里计算类大小，并分配空间，返回生成的对象
    def generateCtorAllocInline(self, ctx):
        if not ctx.isinline: raise '123'

        class_name = ctx.class_name
        method_name = ctx.method_name
        method_cursor = ctx.cursor

        params_arr = self.generateParams(class_name, method_name, method_cursor)
        params = ', '.join(params_arr)

        # 类内类处理
        full_class_name = ctx.full_class_name

        ctx.CP.AP('main', '// %s' % (str(ctx.cursor.location)))
        ctx.CP.AP('main', '// %s' % (ctx.fn_proto_cpp))
        ctx.CP.AP('main', 'if (true) {')
        idx = 0
        argv = []
        prmv = []
        for arg in ctx.cursor.get_arguments():
            idx += 1
            # self.generateParamDeclExpr(ctx, arg, idx)
            argv.append('arg%s' % (idx))
            prme = self.generateParamForInline(arg, idx)
            prmv.append('%s' % (prme))
        args = ', '.join(argv)
        prms = ', '.join(prmv)
        ctx.CP.AP('main', '  auto f = [](%s) {' % (prms))
        if ctx.isabstract:
            ctx.CP.AP('main', '    // new %s(%s);' % (ctx.full_class_name, args))
        else:
            ctx.CP.AP('main', '    new %s(%s);' % (ctx.full_class_name, args))
        ctx.CP.AP('main', '  };')
        ctx.CP.AP('main', '  if (f == nullptr){}')
        ctx.CP.AP('main', '}')

        return

    # 重新生成新的ctor封装，在这里计算类大小，并分配空间，返回生成的对象
    def generateDtorDelete(self, ctx):

        class_name = ctx.class_name
        method_name = ctx.method_name
        method_cursor = ctx.cursor

        full_class_name = ctx.class_name
        # 类内类处理
        if ctx.class_cursor.semantic_parent.kind == clidx.CursorKind.STRUCT_DECL or \
           ctx.class_cursor.semantic_parent.kind == clidx.CursorKind.CLASS_DECL:
            print('ooops', ctx.class_cursor.semantic_parent.kind, ctx.class_cursor.semantic_parent.spelling)
            full_class_name = '%s::%s' % (ctx.class_cursor.semantic_parent.spelling, ctx.class_name)
            # exit(0)

        ctx.CP.AP('ext', '// %s' % (ctx.fn_proto_cpp))
        ctx.CP.AP('ext', 'extern "C"')
        ctx.CP.AP('ext', 'void C%s(void *qthis) {' % (ctx.mangled_name))
        ctx.CP.AP('ext', '  delete (%s*)qthis;' % (ctx.full_class_name))
        ctx.CP.AP('ext', '}')

        return

    # 重新生成新的ctor封装，在这里计算类大小，并分配空间，返回生成的对象
    def generateDtorDeleteInline(self, ctx):
        if not ctx.isinline: raise '123'

        class_name = ctx.class_name
        method_name = ctx.method_name
        method_cursor = ctx.cursor

        full_class_name = ctx.class_name
        # 类内类处理
        if ctx.class_cursor.semantic_parent.kind == clidx.CursorKind.STRUCT_DECL or \
           ctx.class_cursor.semantic_parent.kind == clidx.CursorKind.CLASS_DECL:
            print('ooops', ctx.class_cursor.semantic_parent.kind, ctx.class_cursor.semantic_parent.spelling)
            full_class_name = '%s::%s' % (ctx.class_cursor.semantic_parent.spelling, ctx.class_name)
            # exit(0)

        ctx.CP.AP('main', '// %s' % (ctx.fn_proto_cpp))
        ctx.CP.AP('main', 'if (true) {')
        ctx.CP.AP('main', '  delete ((%s*)0);' % (ctx.full_class_name))
        ctx.CP.AP('main', '}')

        return

    def generateMethod(self, ctx):
        class_name = ctx.class_name
        method_name = ctx.method_name
        method_cursor = ctx.cursor
        cursor = method_cursor

        return_type = cursor.result_type
        return_real_type = self.real_type_name(return_type)
        # if '::' in return_real_type: return
        # if self.check_skip_params(cursor): return

        inner_return = ''
        if cursor.kind == clidx.CursorKind.CONSTRUCTOR or \
           cursor.kind == clidx.CursorKind.DESTRUCTOR:
            pass
        else:
            return_type_name = self.resolve_swig_type_name(class_name, return_type)
            inner_return = 'return' if return_type_name != 'void' else inner_return

        params = self.generateParams(class_name, method_name, method_cursor)
        params = ', '.join(params)
        params = 'void *that, ' + params
        params = params.strip(', ')

        rvref = '&&' in ctx.ret_type_name_cpp and ctx.ret_type_ref
        ret_type_name = ctx.ret_type_name_cpp
        if ctx.ret_type_ref and not rvref: ret_type_name = ctx.ret_type_name_cpp.replace('&', '*')
        mangled_name = method_cursor.mangled_name

        idx = 0
        argv = []
        prmv = []
        if not ctx.static:
            prmv.append("void *qthis")

        # argv.append("ethis")
        for arg in ctx.cursor.get_arguments():
            idx += 1
            # self.generateParamDeclExpr(ctx, arg, idx)
            prme = self.generateParamForDecl(arg, idx)
            prmv.append('%s' % (prme))
            prmc = self.generateParamForCall(arg, idx)
            argv.append(prmc)
            # argv.append('arg%s' % (idx))
        args = ',\n'.join(argv)
        prms = ',\n'.join(prmv)

        ctx.CP.AP('ext', '// %s' % (str(ctx.cursor.location)))
        ctx.CP.AP('ext', '// %s' % (ctx.fn_proto_cpp))
        ctx.CP.AP('ext', '// %s %s' % (ctx.mangled_name, ctx.cursor.displayname))
        ctx.CP.AP('ext', 'extern "C"')
        self.generateReturnDecl(ctx)
        ctx.CP.AP('ext', 'C%s(%s) {' % (ctx.mangled_name, prms))
        if ctx.need_return:
            if ctx.ret_type.kind == clidx.TypeKind.LVALUEREFERENCE:
                ctx.CP.AP('ext', '  auto& ret =')
            else:
                ctx.CP.AP('ext', '  auto ret =')
        if ctx.static:
            ctx.CP.AP('ext', '  %s::%s(%s);' % (ctx.full_class_name, method_name, args))
        else:
            ctx.CP.AP('ext', '  ((%s*)qthis)->%s(%s);' % (ctx.full_class_name, method_name, args))
        self.generateReturnImpl(ctx)
        ctx.CP.AP('ext', '}')

        # 尝试添加正确的#include
        self.generateUseForType(ctx, ctx.ret_type)

        return

    def generateMethodInline(self, ctx):
        if not ctx.isinline: raise '123'

        class_name = ctx.class_name
        method_name = ctx.method_name
        method_cursor = ctx.cursor
        cursor = method_cursor

        return_type = cursor.result_type
        return_real_type = self.real_type_name(return_type)
        # if '::' in return_real_type: return
        # if self.check_skip_params(cursor): return

        inner_return = ''
        if cursor.kind == clidx.CursorKind.CONSTRUCTOR or \
           cursor.kind == clidx.CursorKind.DESTRUCTOR:
            pass
        else:
            return_type_name = self.resolve_swig_type_name(class_name, return_type)
            inner_return = 'return' if return_type_name != 'void' else inner_return

        params = self.generateParams(class_name, method_name, method_cursor)
        params = ', '.join(params)
        params = 'void *that, ' + params
        params = params.strip(', ')

        rvref = '&&' in ctx.ret_type_name_cpp and ctx.ret_type_ref
        ret_type_name = ctx.ret_type_name_cpp
        if ctx.ret_type_ref and not rvref: ret_type_name = ctx.ret_type_name_cpp.replace('&', '*')
        mangled_name = method_cursor.mangled_name

        ctx.CP.AP('main', '// %s' % (str(ctx.cursor.location)))
        ctx.CP.AP('main', '// %s' % (ctx.fn_proto_cpp))
        ctx.CP.AP('main', 'if (true) {')

        idx = 0
        argv = []
        prmv = []
        if not self.gutil.isAbstractClass(ctx.class_cursor):
            prmv.append("%s flythis" % (ctx.full_class_name))
        # argv.append("ethis")
        for arg in ctx.cursor.get_arguments():
            idx += 1
            # self.generateParamDeclExpr(ctx, arg, idx)
            prme = self.generateParamForInline(arg, idx)
            prmv.append('%s' % (prme))
            argv.append('arg%s' % (idx))
        args = ', '.join(argv)
        prms = ', '.join(prmv)
        ctx.CP.AP('main', '  auto f = [](%s) {' % (prms))
        ctx.CP.AP('main', '    ((%s*)0)->%s(%s);' % (ctx.full_class_name, method_name, args))
        if not self.gutil.isAbstractClass(ctx.class_cursor):
            ctx.CP.AP('main', '    flythis.%s(%s);' % (method_name, args))
        ctx.CP.AP('main', '  };')
        ctx.CP.AP('main', '  if (f == nullptr){}')
        ctx.CP.AP('main', '}')
        ctx.CP.AP('main', '// %s %s' % (ctx.mangled_name, ctx.cursor.displayname))

        # 尝试添加正确的#include
        # self.generateUseForType(ctx, ctx.ret_type)

        return

    def generateReturnImpl(self, ctx):
        if not ctx.need_return: return

        ret_type = ctx.cursor.result_type
        ret_type = self.tyconv.TypeToActual(ret_type)

        if ret_type.kind == clidx.TypeKind.VOID:
            if ctx.ctor:
                ctx.CP.AP('ext', '  return ret;')
            else:
                ctx.CP.AP('ext', '  // return void;')
        elif ret_type.kind == clidx.TypeKind.POINTER:
            ctx.CP.AP('ext', '  return (void*)ret;')
        elif ret_type.kind == clidx.TypeKind.RECORD:
            if self.gutil.isDisableCopy(ret_type.get_declaration()):
                ctx.CP.AP('ext', '  return &ret; // return new %s(ret);' % (ret_type.spelling))
            else:
                ctx.CP.AP('ext', '  return new %s(ret); // 5' % (ret_type.spelling))
        elif ret_type.kind == clidx.TypeKind.LVALUEREFERENCE:
            raise '123'
            under_type = ret_type.get_pointee()
            if under_type.kind == clidx.TypeKind.TYPEDEF:
                under_type = under_type.get_declaration().underlying_typedef_type
            if under_type.kind == clidx.TypeKind.UNEXPOSED:
                under_type = under_type.get_declaration().type
            if under_type.kind == clidx.TypeKind.RECORD:
                if self.gutil.isDisableCopy(under_type.get_declaration()):
                    ctx.CP.AP('ext', '  return &ret; // return new %s(ret);' % (under_type.spelling))
                else:
                    ctx.CP.AP('ext', '  return new %s(ret); // 4' % (under_type.spelling))
            else:
                ctx.CP.AP('ext', '  return ret; // 2 %s' % (under_type.kind))
        elif ret_type.kind == clidx.TypeKind.TYPEDEF:
            raise '123'
            under_type = ret_type.get_declaration().underlying_typedef_type
            if under_type.kind == clidx.TypeKind.UNEXPOSED:
                ctx.CP.AP('ext', '  return (void*)(new (decltype(ret))(ret)); // %s' % (ret_type.spelling))
            elif under_type.kind == clidx.TypeKind.RECORD:
                ctx.CP.AP('ext', '  return new %s(ret); // 6' % (ret_type.spelling))
            elif under_type.kind == clidx.TypeKind.FUNCTIONPROTO:
                ctx.CP.AP('ext', '  return (void*)ret;')
            elif under_type.kind == clidx.TypeKind.POINTER:
                ctx.CP.AP('ext', '  return (void*)ret;')
            else:
                ctx.CP.AP('ext', '  return ret; // 1 %s' % (under_type.kind))
        elif ret_type.kind == clidx.TypeKind.UNEXPOSED:
            raise '123'
            ctx.CP.AP('ext', '  return (void*)(new (decltype(ret))(ret)); // %s' % (ret_type.spelling))
        else:
            ctx.CP.AP('ext', '  return ret; // 0 %s' % (ret_type.kind))
        return

    def generateReturnDecl(self, ctx):
        if not ctx.need_return:
            ctx.CP.AP('ext', 'void')
            return

        ret_type = ctx.cursor.result_type
        ret_type = self.tyconv.TypeToActual(ret_type)
        tyname = self.hotfix_class_inner_type(ret_type)

        if ret_type.kind == clidx.TypeKind.VOID:
            if ctx.ctor:
                ctx.CP.AP('ext', '%s*' % (ctx.full_class_name))
            else:
                ctx.CP.AP('ext', 'void')
        elif ret_type.kind == clidx.TypeKind.POINTER:
            ctx.CP.AP('ext', 'void*')
        elif ret_type.kind == clidx.TypeKind.RECORD:
            ctx.CP.AP('ext', '%s*' % (ret_type.spelling))
        elif ret_type.kind == clidx.TypeKind.LVALUEREFERENCE:
            raise '123'
            under_type = self.tyconv.TypeToActual(ret_type)
            if under_type.kind == clidx.TypeKind.RECORD:
                ctx.CP.AP('ext', '%s* /* 1 */' % (under_type.spelling))
            elif under_type.kind == clidx.TypeKind.POINTER:
                ctx.CP.AP('ext', 'void* /* 2 */')
            else:
                ctx.CP.AP('ext', '%s /* 3 */' % (under_type.spelling))
        elif ret_type.kind == clidx.TypeKind.TYPEDEF:
            raise '123'
            under_type = ret_type.get_declaration().underlying_typedef_type
            if under_type.kind == clidx.TypeKind.UNEXPOSED:
                ctx.CP.AP('ext', 'void*  // unexposed2 %s' % (ret_type.spelling))
            elif under_type.kind == clidx.TypeKind.RECORD:
                ctx.CP.AP('ext', '%s*' % (ret_type.spelling))
            elif under_type.kind == clidx.TypeKind.FUNCTIONPROTO:
                ctx.CP.AP('ext', 'void*  // %s' % (ret_type.spelling))
            elif under_type.kind == clidx.TypeKind.POINTER:
                ctx.CP.AP('ext', 'void*')
            else:
                ctx.CP.AP('ext', '%s' % (tyname))
        elif ret_type.kind == clidx.TypeKind.UNEXPOSED:
            raise '123'
            ctx.CP.AP('ext', 'void*  // unexposed %s' % (ret_type.spelling))
        else:
            ctx.CP.AP('ext', '%s' % (ret_type.spelling))
        return

    def generateParamForInline(self, arg, idx):
        aty = arg.type
        tyname = aty.spelling
        tyname = self.hotfix_class_inner_type(aty)

        if aty.kind == clidx.TypeKind.LVALUEREFERENCE:
            return '%s arg%s' % (tyname, idx)
        elif aty.kind == clidx.TypeKind.RVALUEREFERENCE:
            # 引用折叠
            return '%s arg%s' % (aty.spelling.replace('&&', '&&'), idx)
        elif aty.kind == clidx.TypeKind.FUNCTIONNOPROTO:
            return '%s' % (aty.spelling.replace('(*)', '(*arg%s)' % (idx)))
        elif aty.kind == clidx.TypeKind.POINTER and '(*)' in aty.spelling:
            return '%s' % (aty.spelling.replace('(*)', '(*arg%s)' % (idx)))
        elif aty.kind == clidx.TypeKind.INCOMPLETEARRAY:
            return '%s' % (aty.spelling.replace(' [', ' arg%s[' % (idx)))
        elif aty.kind == clidx.TypeKind.CONSTANTARRAY:
            return '%s' % (aty.spelling.replace(' [', ' arg%s[' % (idx)))
        elif aty.kind == clidx.TypeKind.RECORD:
            return '%s arg%s' % (tyname, idx)
        elif aty.kind == clidx.TypeKind.UNEXPOSED:
            return '%s arg%s' % (tyname, idx)
        elif aty.kind == clidx.TypeKind.TYPEDEF:
            under_type = aty.get_declaration().underlying_typedef_type
            if under_type.kind == clidx.TypeKind.RECORD:
                return '%s arg%s' % (tyname, idx)
            else:
                return '%s arg%s' % (tyname, idx)
        else:
            return '%s arg%s' % (tyname, idx)

        # 尝试添加正确的#include

        return

    def generateParamForDecl(self, arg, idx):
        aty = arg.type
        tyname = aty.spelling
        tyname = self.hotfix_class_inner_type(aty)

        if aty.kind == clidx.TypeKind.LVALUEREFERENCE:
            can_type = aty.get_pointee()
            can_tyname = self.hotfix_class_inner_type(can_type)
            return '%s* arg%s' % (can_tyname, idx)
            # return '%s arg%s' % (tyname, idx)
        elif aty.kind == clidx.TypeKind.RVALUEREFERENCE:
            can_type = aty.get_pointee()
            can_tyname = self.hotfix_class_inner_type(can_type)
            return '%s* arg%s' % (can_tyname, idx)
            # 引用折叠
            # return '%s arg%s' % (aty.spelling.replace('&&', '&&'), idx)
        elif aty.kind == clidx.TypeKind.FUNCTIONNOPROTO:
            return '%s' % (aty.spelling.replace('(*)', '(*arg%s)' % (idx)))
        elif aty.kind == clidx.TypeKind.POINTER and '(*)' in aty.spelling:
            return '%s' % (aty.spelling.replace('(*)', '(*arg%s)' % (idx)))
        elif aty.kind == clidx.TypeKind.INCOMPLETEARRAY:
            return '%s' % (aty.spelling.replace(' [', ' arg%s[' % (idx)))
        elif aty.kind == clidx.TypeKind.CONSTANTARRAY:
            return '%s' % (aty.spelling.replace(' [', ' arg%s[' % (idx)))
        elif aty.kind == clidx.TypeKind.RECORD:
            return '%s* arg%s' % (tyname, idx)
        elif aty.kind == clidx.TypeKind.UNEXPOSED:
            return '%s* arg%s' % (tyname, idx)
        elif aty.kind == clidx.TypeKind.TYPEDEF:
            under_type = aty.get_declaration().underlying_typedef_type
            if under_type.kind == clidx.TypeKind.RECORD:
                return '%s* arg%s' % (tyname, idx)
            else:
                return '%s arg%s' % (tyname, idx)
        else:
            return '%s arg%s' % (tyname, idx)

        # 尝试添加正确的#include

        return

    def generateParamForCall(self, arg, idx):
        aty = arg.type
        tyname = aty.spelling
        tyname = self.hotfix_class_inner_type(aty)

        if aty.kind == clidx.TypeKind.LVALUEREFERENCE:
            can_type = aty.get_pointee()
            can_tyname = self.hotfix_class_inner_type(can_type)
            return '*((%s*)arg%s)' % (can_tyname, idx)
            # return '%s arg%s' % (tyname, idx)
        elif aty.kind == clidx.TypeKind.RVALUEREFERENCE:
            can_type = aty.get_pointee()
            can_tyname = self.hotfix_class_inner_type(can_type)
            return '*((%s*)arg%s)' % (can_tyname, idx)
            # 引用折叠
            # return '%s arg%s' % (aty.spelling.replace('&&', '&&'), idx)
        elif aty.kind == clidx.TypeKind.FUNCTIONNOPROTO:
            return 'arg%s' % (idx)
            # return '%s' % (aty.spelling.replace('(*)', '(*arg%s)' % (idx)))
        elif aty.kind == clidx.TypeKind.POINTER and '(*)' in aty.spelling:
            return 'arg%s' % (idx)
            # return '%s' % (aty.spelling.replace('(*)', '(*arg%s)' % (idx)))
        elif aty.kind == clidx.TypeKind.INCOMPLETEARRAY:
            return 'arg%s' % (idx)
            # return '%s' % (aty.spelling.replace(' [', ' arg%s[' % (idx)))
        elif aty.kind == clidx.TypeKind.CONSTANTARRAY:
            return 'arg%s' % (idx)
            # return '%s' % (aty.spelling.replace(' [', ' arg%s[' % (idx)))
        elif aty.kind == clidx.TypeKind.RECORD:
            return '*((%s*)arg%s)' % (tyname, idx)
            # return '%s* arg%s' % (tyname, idx)
        elif aty.kind == clidx.TypeKind.UNEXPOSED:
            return '*((%s*)arg%s)' % (tyname, idx)
            # return '%s* arg%s' % (tyname, idx)
        elif aty.kind == clidx.TypeKind.TYPEDEF:
            under_type = aty.get_declaration().underlying_typedef_type
            if under_type.kind == clidx.TypeKind.RECORD:
                return '*((%s*)arg%s)' % (tyname, idx)
                # return '%s* arg%s' % (tyname, idx)
            else:
                return 'arg%s' % (idx)
        else:
            return 'arg%s' % (idx)

        # 尝试添加正确的#include

        return

    # @return []
    def generateParamsRaw(self, class_name, method_name, method_cursor):
        argv = []
        for arg in method_cursor.get_arguments():
            argelem = "%s %s" % (arg.type.spelling, arg.displayname)
            argv.append(argelem)
            if '<' in arg.type.spelling:
                xdef = arg.type.get_declaration()
                tic = self.gutil.isTempInstClass(xdef)
                if tic is not None:
                    class_cursor = self.get_instantiated_class(xdef)
                    mths = self.gutil.get_inst_methods(class_cursor, xdef)
        return argv

    # @return []
    def generateParams(self, class_name, method_name, method_cursor):
        idx = 0
        argv = []

        for arg in method_cursor.get_arguments():
            idx += 1
            # print('%s, %s, ty:%s, kindty:%s' % (method_name, arg.displayname, arg.type.spelling, arg.kind))
            # print('arg type kind: %s, %s' % (arg.type.kind, arg.type.get_declaration()))
            type_name = self.resolve_swig_type_name(class_name, arg.type)

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
            type_name = self.resolve_swig_type_name(class_name, arg.type)

            arg_name = 'arg%s' % idx if arg.displayname == '' else arg.displayname
            argelem = "%s" % (arg_name)
            argv.append(argelem)

        return argv

    # 对于C封装来说的
    def methodNeedReturn(self, ctx):
        ret_type = ctx.cursor.result_type

        if ret_type.kind == clidx.TypeKind.VOID:
            if ctx.ctor: return True
            else: return False
        else:
            return True
        return False

    def methodHasReturn(self, ctx):
        method_cursor = cursor = ctx.cursor
        class_name = ctx.class_name

        return_type = cursor.result_type

        return_type_name = return_type.spelling
        if ctx.ctor or ctx.dtor: pass
        else:
            return_type_name = self.resolve_swig_type_name(class_name, return_type)

        has_return = True
        if return_type.kind == clidx.TypeKind.VOID:
            has_return = False
        # if return_type_name == 'void': has_return = False
        # if cursor.spelling == 'buttons':
        #     print(666, has_return, return_type_name, cursor.spelling, return_type.kind, cursor.semantic_parent.spelling)
        #     exit(0)
        # if return_type_name.count('<') >= 2:
        #     has_return = False
        # elif return_type_name.count('<') == 1:
        #     # print(556, return_type_name, ctx.fn_proto_cpp)
        #     xdef = return_type.get_declaration()
        #     tic = self.gutil.isTempInstClass(xdef)
        #     if tic is not None:
        #         class_cursor = self.get_instantiated_class(xdef)
        #         mths = self.gutil.get_inst_methods(class_cursor, xdef)
        #     has_return = False

        # if '::' in return_type_name: has_return = False
        # if "QStringList" in return_type_name: has_return = False
        # if "QObjectList" in return_type_name: has_return = False
        # if 'QAbstract' in return_type_name: has_return = False
        # if 'QMetaObject' in return_type_name: has_return = False
        # if 'QOpenGL' in return_type_name: has_return = False
        # if 'QGraphics' in return_type_name: has_return = False
        # if 'QPlatform' in return_type_name: has_return = False
        # if 'QFunctionPointer' in return_type_name: has_return = False
        # if 'QTextEngine' in return_type_name: has_return = False
        # if 'QTextDocumentPrivate' in return_type_name: has_return = False
        # if 'QJson' in return_type_name: has_return = False
        # if 'QStringRef' in return_type_name: has_return = False

        # if 'internalPointer' in method_cursor.spelling: has_return = False
        # if 'rwidth' in method_cursor.spelling: has_return = False
        # if 'rheight' in method_cursor.spelling: has_return = False
        # if 'utf16' == method_cursor.spelling: has_return = False
        # if 'x' == method_cursor.spelling: has_return = False
        # if 'rx' == method_cursor.spelling: has_return = False
        # if 'y' == method_cursor.spelling: has_return = False
        # if 'ry' == method_cursor.spelling: has_return = False
        # if class_name == 'QGenericArgument' and method_cursor.spelling == 'data': has_return = False
        # if class_name == 'QSharedMemory' and method_cursor.spelling == 'constData': has_return = False
        # if class_name == 'QSharedMemory' and method_cursor.spelling == 'data': has_return = False
        # if class_name == 'QVariant' and method_cursor.spelling == 'constData': has_return = False
        # if class_name == 'QVariant' and method_cursor.spelling == 'data': has_return = False
        # if class_name == 'QThreadStorageData' and method_cursor.spelling == 'set': has_return = False
        # if class_name == 'QThreadStorageData' and method_cursor.spelling == 'get': has_return = False
        # if class_name == 'QChar' and method_cursor.spelling == 'unicode': has_return = False

        return has_return

    # cursor, 当前位置，可以提供能确定所在文件的值
    # cty, 如参数类型，或者返回值类型
    def generateUseForType(self, ctx, cty):
        if cty.kind == clidx.TypeKind.LVALUEREFERENCE \
           or cty.kind == clidx.TypeKind.LVALUEREFERENCE:
            cty = cty.get_pointee()
        xdef = cty.get_declaration()
        if xdef is not None:
            cloc = ctx.cursor.location.file
            xloc = xdef.location.file

            if cty.spelling.count('<') >= 2:
                pass
            elif cty.spelling.count('<') == 1 and cty.spelling.startswith('Q'):
                tic = self.gutil.isTempInstClass(xdef)
                if tic is not None:
                    self.gutil.ticlasses[cty.spelling] = xdef
                    class_cursor = self.get_instantiated_class(xdef)
                    mths = self.gutil.get_inst_methods(class_cursor, xdef)

            if '<' in cty.spelling and cty.spelling.startswith('Q'):
                chname = cloc.name.split('/')[-1]
                xhname = xloc.name.split('/')[-1]
                thname = cty.spelling.split('<')[0].lower() + '.h'
                if chname != thname and thname[0] == 'q':
                    ctx.CP.APU('header', '#include <%s>' % (thname))

            if xloc is not None and xloc.name != cloc.name:
                chname = cloc.name.split('/')[-1]
                xhname = xloc.name.split('/')[-1]
                if xhname[0] == 'q':
                    ctx.CP.APU('header', '#include <%s>' % (xhname))
        return

    def is_conflict_method_name(self, method_name):
        return False
        if method_name in ['match', 'type', 'move']:  # , 'select']:
            return True
        return False

    # @return True | False
    def check_skip_params(self, cursor):
        method_name = cursor.spelling
        for arg in cursor.get_arguments():
            if self.check_skip_param(arg, method_name) is True: return True
        return False

    def check_skip_param(self, arg, method_name):
        if True:
            type_name = arg.type.spelling
            type_name_segs = type_name.split(' ')
            if 'const' in type_name_segs: type_name_segs.remove('const')
            if '*' in type_name_segs: type_name_segs.remove('*')
            if '&' in type_name_segs: type_name_segs.remove('&')
            type_name = type_name_segs[0]

            # shitfix
            atydecl = arg.type.get_declaration()
            if atydecl.access_specifier == clidx.AccessSpecifier.PRIVATE \
               or atydecl.access_specifier == clidx.AccessSpecifier.PROTECTED:
                # print(87987, arg.type.kind, atydecl.access_specifier, atydecl.spelling)
                return True

        return False

    # @return True | False
    def check_skip_method(self, cursor):
        # shitfix
        if cursor.mangled_name == '_ZN14QSignalBlockerC1EOS_': return True
        if cursor.mangled_name == '_ZN15QAnimationGroupC1EP7QObject': return True  # abstract
        if cursor.mangled_name == '_ZN17QAccessibleObjectC1EP7QObject': return True  # abstract

        if True: return self.gfilter.skipMethod(cursor)
        return False

    def check_skip_class(self, class_cursor):
        if True: return self.gfilter.skipClass(class_cursor)
        return False

    # 类似，QSurfaceFormat::FormatOptions
    def hotfix_class_inner_type(self, cty):
        tyname = cty.spelling
        xdef = cty.get_declaration()
        if cty.kind == clidx.TypeKind.LVALUEREFERENCE \
           or cty.kind == clidx.TypeKind.POINTER:
            xdef = cty.get_pointee().get_declaration()
        if xdef is not None and xdef.kind == clidx.TypeKind.TYPEDEF:
            xdef = xdef.type.get_declaration()

        def removeQuality(tyname):
            lst = tyname.replace('*', ' * ').split()
            nlst = []
            for e in lst:
                if e not in ['const', '*', '&']:
                    nlst.append(e)
            return ' '.join(nlst)

        if xdef is not None and xdef.semantic_parent is not None:
            pdef = xdef.semantic_parent
            if pdef.kind == clidx.CursorKind.CLASS_DECL \
               or pdef.kind == clidx.CursorKind.STRUCT_DECL:
                if '::' not in tyname and tyname[0] != 'Q':
                    ctyname = removeQuality(tyname)
                    ntyname = tyname.replace(ctyname, '%s::%s' % (pdef.spelling, ctyname))
                    print(666890, tyname, '=>', ntyname)
                    tyname = ntyname
        return tyname

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
            if mod not in ['core', 'gui', 'widgets', 'network', 'dbus', 'qml', 'quick']:
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

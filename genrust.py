# encoding: utf8

import os
import logging
import time

import clang
import clang.cindex
import clang.cindex as clidx

from genutil import *
from typeconv import TypeConv, TypeConvForRust
from genbase import GenerateBase

# TODO
# 静态方法与动态方法同名
# 参数默认值
# 内联方法
# enum类型成员
# qt 全局函数
# C++类的继承方法
# 继承依赖继承链，而中间可能有QAbstractxxx类，需要处理。
# 集合参数或返回值的转换，像Vec<T> <=> QList<T>, 或者Vec<T> <=> T **
# qt模板类型的封装实现
# 代码整理, GenContext -- OK
# 生成个简单文档？然后生成文档。


class GenMethodContext(object):
    def __init__(self, cursor, class_cursor):
        self.tyconv = TypeConvForRust()

        self.ctysz = max(32, class_cursor.type.get_size())  # 可能这个get_size()的值不准确啊。
        self.class_cursor = class_cursor
        self.class_name = class_cursor.spelling
        self.cursor = cursor
        self.method_name = cursor.spelling
        self.method_name_rewrite = self.method_name
        self.mangled_name = cursor.mangled_name

        self.ctor = cursor.kind == clidx.CursorKind.CONSTRUCTOR
        self.dtor = cursor.kind == clidx.CursorKind.DESTRUCTOR

        self.static = cursor.is_static_method()
        self.has_return = True
        self.ret_type = cursor.result_type
        self.ret_type_name_cpp = self.ret_type.spelling
        self.ret_type_name_rs = ''
        self.ret_type_name_ext = ''

        self.static_str = 'static' if self.static else ''
        self.static_suffix = '_s' if self.static else ''
        self.static_self_struct = '' if self.static else '&mut self, '
        self.static_self_trait = '' if self.static else ', rsthis: &mut %s' % (self.class_name)
        self.static_self_call = '' if self.static else 'self'

        self.params_cpp = ''
        self.params_rs = ''
        self.params_call = ''
        self.params_ext = ''

        self.unique_methods = {}
        self.struct_proto = '%s::%s%s' % (self.class_name, self.method_name, self.static_suffix)
        self.trait_proto = ''  # '%s::%s(%s)' % (class_name, method_name, trait_params)

        self.fn_proto_cpp = ''

        # inherit
        self.base_class = None
        self.base_class_name = ''
        self.has_base = False

        # aux
        self.tymap = None
        # simple init

        self.CP = None
        return


class GenerateForRust(GenerateBase):
    def __init__(self):
        super(GenerateForRust, self).__init__()

        self.modrss = {}  # mod => CodePaper
        #self.cp_modrs = CodePaper()  # 可能的name: main
        #self.cp_modrs.addPoint('main')
        #self.MP = self.cp_modrs

        self.class_blocks = ['header', 'main', 'use', 'ext', 'body']
        # self.cp_clsrs = CodePaper()  # 可能中间reset。可能的name: header, main, use, ext, body
        # self.CP = self.cp_clsrs

        self.qclses = {}  # class name => True
        self.tyconv = TypeConvForRust()
        self.traits = {}  # traits proto => True
        self.implmthods = {}  # method proto => True
        return

    def generateHeader(self, module):
        code = ''
        code += "#![feature(libc)]\n"
        code += "#![feature(core)]\n"
        code += "#![feature(collections)]\n"
        code += "extern crate libc;\n"
        code += "use self::libc::*;\n"

        code += "\n"
        return code

    def generateFooter(self, module):
        return ''

    def initCodePaperForClass(self):
        cp_clsrs = CodePaper()
        for blk in self.class_blocks:
            cp_clsrs.addPoint(blk)
            cp_clsrs.append(blk, "// %s block begin" % (blk))
        return cp_clsrs

    def genpass_init_code_paper(self):
        for key in self.gctx.codes:
            CP = self.gctx.codes[key]
            mod = self.gctx.get_decl_mod_by_path(key)
            code_file = self.gctx.get_code_file_by_path(key)
            CP.AP('header', '// auto generated, do not modify.')
            CP.AP('header', '// created: ' + time.ctime())
            CP.AP('header', '// src-file: ' + key)
            CP.AP('header', '// dst-file: /src/%s/%s.rs' % (mod, code_file))
            CP.AP('header', '//\n')

            for blk in self.class_blocks:
                CP.addPoint(blk)
                CP.append(blk, "// %s block begin =>" % (blk))
        return

    def genpass_code_header(self):
        for key in self.gctx.codes:
            CP = self.gctx.codes[key]
            CP.AP('header', self.generateHeader(''))
            CP.AP('ext', "#[link(name = \"Qt5Core\")]")
            CP.AP('ext', "#[link(name = \"Qt5Gui\")]")
            CP.AP('ext', "#[link(name = \"Qt5Widgets\")]\n")
            CP.AP('ext', "extern {")

        return

    def genpass_code_endian(self):
        for key in self.gctx.codes:
            CP = self.gctx.codes[key]
            for blk in self.class_blocks:
                if blk == 'ext':
                    CP.append(blk, "} // <= %s block end\n" % (blk))
                else:
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
            MP.APU('main', "pub mod %s;" % (code_file))
            MP.APU('main', "pub use self::%s::%s;\n" % (code_file, class_name))
        return

    def genpass(self):
        self.genpass_init_code_paper()
        self.genpass_code_header()

        self.genpass_class_type()

        print('gen classes...')
        self.genpass_classes()

        print('gen code endian...')
        self.genpass_code_endian()

        print('gen class mod define...')
        self.genpass_class_modef()

        print('gen files...')
        self.genpass_write_codes()
        return

    def genpass_class_type(self):
        for key in self.gctx.classes:
            cursor = self.gctx.classes[key]
            if self.check_skip_class(cursor): continue
            self.genpass_class_type_impl(cursor)
        return

    def genpass_class_type_impl(self, cursor):
        class_name = cursor.displayname
        decl_file = self.gctx.get_decl_file(cursor)
        CP = self.gctx.codes[decl_file]
        ctysz = cursor.type.get_size()

        CP.AP('body', "// class sizeof(%s)=%s" % (class_name, ctysz))
        # generate struct of class
        CP.AP('body', "pub struct %s {" % (class_name))
        CP.AP('body', "  pub qclsinst: *mut c_void,")
        CP.AP('body', "}\n")

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

    # def generateClasses(self, module, class_decls):
    #     for elems in class_decls:
    #         class_name, cs, methods, base_class = elems
    #         self.qclses[class_name] = True

    #     for elems in class_decls:
    #         class_name, cs, methods, base_class = elems
    #         self.CP = self.initCodePaperForClass()
    #         self.CP.AP('header', self.generateHeader(module))
    #         self.CP.AP('ext', "#[link(name = \"Qt5Core\")]\n")
    #         self.CP.AP('ext', "#[link(name = \"Qt5Gui\")]\n")
    #         self.CP.AP('ext', "#[link(name = \"Qt5Widgets\")]\n")
    #         self.CP.AP('ext', "extern {\n")

    #         self.generateClass(class_name, cs, methods, base_class)
    #         # tcode = tcode + self.generateFooter(module)
    #         # self.write_code(module, class_name.lower(), tcode)
    #         self.CP.AP('ext', "}\n\n")
    #         self.CP.AP('use', "\n")

    #         self.write_code(module, class_name.lower(), self.CP.exportCode(self.class_blocks))

    #         self.MP.AP('main', "mod %s;\n" % (class_name.lower()))
    #         self.MP.AP('main', "pub use self::%s::%s;\n\n" % (class_name.lower(), class_name))

    #     self.write_modrs(module, self.MP.exportCode(['main']))
    #     return

    def generateClass(self, class_name, class_cursor, methods, base_class):

        CP = self.gctx.getCodePager(class_cursor)

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

        dupremove = self.dedup_return_const_diff_method(methods)
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

        if ctx.ctor: ctx.method_name_rewrite = 'New%s' % (ctx.method_name)
        if ctx.dtor: ctx.method_name_rewrite = 'Free%s' % (ctx.method_name[1:])
        if self.is_conflict_method_name(ctx.method_name):
            ctx.method_name_rewrite = ctx.method_name + '_'
        if ctx.static:
            ctx.method_name_rewrite = ctx.method_name + ctx.static_suffix

        class_name = ctx.class_name
        method_name = ctx.method_name

        ctx.ret_type_name_rs = self.tyconv.Type2RustRet(ctx.ret_type, method_cursor)
        ctx.ret_type_name_ext = self.tyconv.TypeCXX2RustExtern(ctx.ret_type, method_cursor)

        raw_params_array = self.generateParamsRaw(class_name, method_name, method_cursor)
        raw_params = ', '.join(raw_params_array)

        trait_params_array = self.generateParamsForTrait(class_name, method_name, method_cursor, ctx)
        trait_params = ', '.join(trait_params_array)

        call_params_array = self.generateParamsForCall(class_name, method_name, method_cursor)
        call_params = ', '.join(call_params_array)
        if not ctx.static and not ctx.ctor: call_params = ('rsthis.qclsinst, ' + call_params).strip(' ,')

        extargs_array = self.generateParamsForExtern(class_name, method_name, method_cursor, ctx)
        extargs = ', '.join(extargs_array)
        if not ctx.static: extargs = ('qthis: *mut c_void, ' + extargs).strip(' ,')

        ctx.params_cpp = raw_params
        ctx.params_rs = trait_params
        ctx.params_call = call_params
        ctx.params_ext = extargs

        ctx.trait_proto = '%s::%s(%s)' % (class_name, method_name, trait_params)
        ctx.fn_proto_cpp = "  // proto: %s %s %s::%s(%s);" % \
                           (ctx.static_str, ctx.ret_type_name_cpp, ctx.class_name, ctx.method_name, ctx.params_cpp)
        ctx.has_return = self.methodHasReturn(ctx)

        # base class
        ctx.base_class = base_class
        ctx.base_class_name = base_class.spelling if base_class is not None else ''
        ctx.has_base = True if base_class is not None else False

        # aux
        ctx.tymap = TypeConvForRust.tymap

        return ctx

    def generateMethod(self, ctx):
        cursor = ctx.cursor

        return_type = cursor.result_type
        return_real_type = self.real_type_name(return_type)
        if '::' in return_real_type: return
        if self.check_skip_params(cursor): return

        static_suffix = ctx.static_suffix

        # method impl
        impl_method_proto = ctx.struct_proto
        if impl_method_proto not in self.implmthods:
            self.implmthods[impl_method_proto] = True
            if ctx.ctor is True: self.generateImplStructCtor(ctx)
            else: self.generateImplStructMethod(ctx)

        uniq_method_name = cursor.spelling + static_suffix
        if ctx.unique_methods[uniq_method_name] is True:
            ctx.unique_methods[uniq_method_name] = False
            self.generateMethodDeclTrait(ctx)

        ### trait impl
        if ctx.trait_proto not in self.traits:
            self.traits[ctx.trait_proto] = True
            if ctx.ctor is True: self.generateImplTraitCtor(ctx)
            else: self.generateImplTraitMethod(ctx)

        # extern
        ctx.CP.AP('ext', ctx.fn_proto_cpp)
        self.generateDeclForFFIExt(ctx)

        return

    def generateImplStructCtor(self, ctx):
        class_name = ctx.class_name
        method_name = ctx.method_name_rewrite

        ctx.CP.AP('body', ctx.fn_proto_cpp)
        ctx.CP.AP('body', "impl /*struct*/ %s {" % (class_name))
        ctx.CP.AP('body', "  pub fn %s<T: %s_%s>(value: T) -> %s {"
                   % (method_name, class_name, method_name, class_name))
        ctx.CP.AP('body', "    let rsthis = value.%s();" % (method_name))
        ctx.CP.AP('body', "    return rsthis;")
        ctx.CP.AP('body', "    // return 1;")
        ctx.CP.AP('body', "  }")
        ctx.CP.AP('body', "}\n")
        return

    def generateImplStructMethod(self, ctx):
        class_name = ctx.class_name
        method_name = ctx.method_name_rewrite
        self_code_proto = ctx.static_self_struct
        self_code_call = ctx.static_self_call

        ctx.CP.AP('body', ctx.fn_proto_cpp)
        ctx.CP.AP('body', "impl /*struct*/ %s {" % (class_name))
        ctx.CP.AP('body', "  pub fn %s<RetType, T: %s_%s<RetType>>(%s overload_args: T) -> RetType {"
                   % (method_name, class_name, method_name, self_code_proto))
        ctx.CP.AP('body', "    return overload_args.%s(%s);" % (method_name, self_code_call))
        ctx.CP.AP('body', "    // return 1;")
        ctx.CP.AP('body', "  }")
        ctx.CP.AP('body', "}\n")
        return

    def generateImplTraitCtor(self, ctx):
        method_cursor = ctx.cursor
        mangled_name = ctx.mangled_name
        class_name = ctx.class_name
        method_name = ctx.method_name_rewrite
        trait_params = ctx.params_rs
        call_params = ctx.params_call

        ctx.CP.AP('body', ctx.fn_proto_cpp)
        ctx.CP.AP('body', "impl<'a> /*trait*/ %s_%s for (%s) {" % (class_name, method_name, trait_params))
        ctx.CP.AP('body', "  fn %s(self) -> %s {" % (method_name, class_name))
        ctx.CP.AP('body', "    let qthis: *mut c_void = unsafe{calloc(1, %s)};" % (ctx.ctysz))
        ctx.CP.AP('body', "    // unsafe{%s()};" % (mangled_name))
        self.generateArgConvExprs(class_name, method_name, method_cursor, ctx)
        if len(call_params) == 0:
            ctx.CP.AP('body', "    unsafe {%s(qthis%s)};" % (mangled_name, call_params))
        else:
            ctx.CP.AP('body', "    unsafe {%s(qthis, %s)};" % (mangled_name, call_params))
        ctx.CP.AP('body', "    let rsthis = %s{qclsinst: qthis};" % (class_name))
        ctx.CP.AP('body', "    return rsthis;")
        ctx.CP.AP('body', "    // return 1;")
        ctx.CP.AP('body', "  }")
        ctx.CP.AP('body', "}\n")

        return

    def generateImplTraitMethod(self, ctx):
        class_name = ctx.class_name
        method_cursor = cursor = ctx.cursor
        method_name = ctx.method_name_rewrite

        has_return = ctx.has_return
        return_piece_code_return = ''
        return_type_name_rs = '()'
        if has_return:
            return_type_name_rs = ctx.ret_type_name_rs
            # print(890, cursor.result_type.spelling, '=>', return_type_name_rs)
            return_piece_code_return = 'let mut ret ='

        self_code_proto = ctx.static_self_trait
        trait_params = ctx.params_rs
        call_params = ctx.params_call

        mangled_name = ctx.mangled_name
        ctx.CP.AP('body', ctx.fn_proto_cpp)
        ctx.CP.AP('body', "impl<'a> /*trait*/ %s_%s<%s> for (%s) {" %
                   (class_name, method_name, return_type_name_rs, trait_params))
        ctx.CP.AP('body', "  fn %s(self %s) -> %s {" %
                   (method_name, self_code_proto, return_type_name_rs))
        ctx.CP.AP('body', "    // let qthis: *mut c_void = unsafe{calloc(1, %s)};" % (ctx.ctysz))
        ctx.CP.AP('body', "    // unsafe{%s()};" % (mangled_name))
        self.generateArgConvExprs(class_name, method_name, method_cursor, ctx)
        ctx.CP.AP('body', "    %s unsafe {%s(%s)};" % (return_piece_code_return, mangled_name, call_params))

        def iscvoidstar(tyname): return ' c_void' in tyname and '*' in tyname
        def isrstar(tyname): return '*' in tyname

        # return expr post process
        # TODO 还有一种值返回的情况要处理，值返回的情况需要先创建一个空对象
        return_type_name_ext = ctx.ret_type_name_ext
        return_type_name_rs = ctx.ret_type_name_rs
        if return_type_name_rs == 'String' and 'char' in return_type_name_ext:
            if has_return: ctx.CP.AP('body', "    let slen = unsafe {strlen(ret as *const i8)} as usize;")
            if has_return: ctx.CP.AP('body', "    return unsafe{String::from_raw_parts(ret as *mut u8, slen, slen+1)};")
        # elif return_type_name_ext == '*mut c_void' or return_type_name_ext == '*const c_void':  # no const now
        elif iscvoidstar(return_type_name_ext) and not isrstar(return_type_name_rs):
            # 应该是返回一个qt class对象，由于无法返回&mut类型的对象
            if has_return: ctx.CP.AP('body', "    let mut ret1 = %s{qclsinst: ret};" % (return_type_name_rs))
            if has_return: ctx.CP.AP('body', "    return ret1;")
        else:
            if has_return: ctx.CP.AP('body', "    return ret as %s;" % (return_type_name_rs))

        ctx.CP.AP('body', "    // return 1;")
        ctx.CP.AP('body', "  }")
        ctx.CP.AP('body', "}\n")

        # case for return qt object
        if has_return:
            self.generateUseForRust(ctx, ctx.ret_type, ctx.cursor)
            # return_type_name = ctx.ret_type_name_rs
            # if self.is_qt_class(return_type_name):
            #     seg = self.get_qt_class(return_type_name)
            #     if seg != class_name and class_name:
            #         ctx.CP.APU('use', "use super::%s::%s;\n" % (seg.lower(), seg))

        return

    def generateMethodDeclTrait(self, ctx):
        class_name = ctx.class_name
        method_name = ctx.method_name_rewrite

        self_code_proto = ctx.static_self_trait

        ### trait
        if ctx.ctor is True:
            ctx.CP.AP('body', "pub trait %s_%s {" % (class_name, method_name))
            ctx.CP.AP('body', "  fn %s(self) -> %s;" % (method_name, class_name))
        else:
            ctx.CP.AP('body', "pub trait %s_%s<RetType> {" % (class_name, method_name))
            ctx.CP.AP('body', "  fn %s(self %s) -> RetType;" %
                       (method_name, self_code_proto))
        ctx.CP.AP('body', "}\n")
        return

    def generateArgConvExprs(self, class_name, method_name, method_cursor, ctx):
        argc = 0
        for arg in method_cursor.get_arguments(): argc += 1

        def isvec(tyname): return 'Vec<' in tyname
        def isrstr(tyname): return 'String' in tyname.split(' ')

        for idx, (arg) in enumerate(method_cursor.get_arguments()):
            srctype = self.tyconv.TypeCXX2Rust(arg.type, arg)
            astype = self.tyconv.TypeCXX2RustExtern(arg.type, arg)
            astype = ' as %s' % (astype)
            asptr = ''
            if self.tyconv.IsPointer(arg.type) and self.tyconv.IsCharType(arg.type.spelling):
                asptr = '.as_ptr()'
            elif isvec(srctype): asptr = '.as_ptr()'
            elif isrstr(srctype): asptr = '.as_ptr()'

            qclsinst = ''
            can_name = self.tyconv.TypeCanName(arg.type)
            if self.is_qt_class(can_name): qclsinst = '.qclsinst'
            if argc == 1:  # fix shit rust tuple index
                ctx.CP.AP('body', "    let arg%s = self%s%s %s;" % (idx, qclsinst, asptr, astype))
            else:
                ctx.CP.AP('body', "    let arg%s = self.%s%s%s %s;" % (idx, idx, qclsinst, asptr, astype))
        return

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
            if arg.type.kind == clang.cindex.TypeKind.POINTER and "(*)" in type_name:
                argelem = "%s" % (type_name.replace('(*)', '(*%s)' % arg_name))
            else:
                argelem = "%s %s" % (type_name, arg_name)
            argv.append(argelem)

        return argv

    # @return []
    def generateParamsRaw(self, class_name, method_name, method_cursor):
        argv = []
        for arg in method_cursor.get_arguments():
            argelem = "%s %s" % (arg.type.spelling, arg.displayname)
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

            type_name_extern = self.tyconv.TypeCXX2RustExtern(arg.type, arg)
            arg_name = 'arg%s' % idx if arg.displayname == '' else arg.displayname
            argelem = "arg%s" % (idx - 1)
            argv.append(argelem)

        return argv

    # @return []
    def generateParamsForTrait(self, class_name, method_name, method_cursor, ctx):
        idx = 0
        argv = []

        for arg in method_cursor.get_arguments():
            idx += 1
            # print('%s, %s, ty:%s, kindty:%s' % (method_name, arg.displayname, arg.type.spelling, arg.kind))
            # print('arg type kind: %s, %s' % (arg.type.kind, arg.type.get_declaration()))
            # param_line2 = self.restore_param_by_token(arg)
            # print(param_line2)

            if self.check_skip_param(arg, method_name) is False:
                self.generateUseForRust(ctx, arg.type, arg)

            # type_name = self.resolve_swig_type_name(class_name, arg.type)
            # type_name2 = self.hotfix_typename_ifenum_asint(class_name, arg, arg.type)
            # type_name = type_name2 if type_name2 is not None else type_name
            type_name = self.tyconv.TypeCXX2Rust(arg.type, arg, method_name)
            if type_name.startswith('&'): type_name = type_name.replace('&', "&'a ")
            # if self.is_qt_class(type_name) and self.check_skip_param(arg, method_name) is False:
            #     seg = self.get_qt_class(type_name)
            #     if seg != class_name:
            #         ctx.CP.APU('use', "use super::%s::%s;\n" % (seg.lower(), seg))

            arg_name = 'arg%s' % idx if arg.displayname == '' else arg.displayname
            argelem = "%s" % (type_name)
            argv.append(argelem)

        return argv

    # @return []
    def generateParamsForExtern(self, class_name, method_name, method_cursor, ctx):
        idx = 0
        argv = []

        if method_cursor.kind == clang.cindex.CursorKind.CONSTRUCTOR:
            # argv.append('qthis: *mut c_void')
            pass

        for arg in method_cursor.get_arguments():
            idx += 1
            # print('%s, %s, ty:%s, kindty:%s' % (method_name, arg.displayname, arg.type.spelling, arg.kind))
            # print('arg type kind: %s, %s' % (arg.type.kind, arg.type.get_declaration()))
            # param_line2 = self.restore_param_by_token(arg)
            # print(param_line2)

            if self.check_skip_param(arg, method_name) is False:
                self.generateUseForRust(ctx, arg.type, arg)
            # type_name = self.resolve_swig_type_name(class_name, arg.type)
            # type_name2 = self.hotfix_typename_ifenum_asint(class_name, arg, arg.type)
            # type_name = type_name2 if type_name2 is not None else type_name
            type_name = self.tyconv.TypeCXX2RustExtern(arg.type, arg)
            # if self.is_qt_class(type_name) and self.check_skip_param(arg, method_name) is False:
            #    seg = self.get_qt_class(type_name)
            #    if seg != class_name:
            #        ctx.CP.APU('use', "use super::%s::%s;\n" % (seg.lower(), seg))

            arg_name = 'arg%s' % idx if arg.displayname == '' else arg.displayname
            argelem = "arg%s: %s" % (idx-1, type_name)
            argv.append(argelem)

        return argv

    def generateReturnForImplStruct(self, class_name, method_cursor, ctx):
        cursor = ctx.cursor

        return_type = cursor.result_type
        return_real_type = self.real_type_name(return_type)

        return_type_name = return_type.spelling
        if cursor.kind == clang.cindex.CursorKind.CONSTRUCTOR or \
           cursor.kind == clang.cindex.CursorKind.DESTRUCTOR:
            pass
        else:
            return_type_name = self.resolve_swig_type_name(class_name, return_type)
            return_type_name2 = self.hotfix_typename_ifenum_asint(class_name, method_cursor, return_type)
            return_type_name = return_type_name2 if return_type_name2 is not None else return_type_name

        has_return = ctx.has_return

        return has_return, return_type_name

    def generateDeclForFFIExt(self, ctx):
        cursor = ctx.cursor
        has_return = ctx.has_return
        # calc ext type name
        return_type_name = self.tyconv.TypeCXX2RustExtern(ctx.ret_type, cursor)

        mangled_name = ctx.mangled_name
        return_piece_proto = ''
        if cursor.result_type.kind != clang.cindex.TypeKind.VOID and has_return:
            return_piece_proto = ' -> %s' % (return_type_name)
        extargs = ctx.params_ext
        ctx.CP.AP('ext', "  fn %s(%s)%s;" % (mangled_name, extargs, return_piece_proto))

        return has_return, return_type_name

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

    def generateUseForRust(self, ctx, aty, cursor):
        class_name = ctx.class_name
        # type_name = self.resolve_swig_type_name(class_name, arg.type)
        # type_name2 = self.hotfix_typename_ifenum_asint(class_name, arg, arg.type)
        # type_name = type_name2 if type_name2 is not None else type_name
        type_name = self.tyconv.TypeCXX2Rust(aty, cursor)
        if type_name.startswith('&'): type_name = type_name.replace('&', "&'a ")
        if self.is_qt_class(type_name):
            seg = self.get_qt_class(type_name)
            # 不但不能是当前类，并且也不能是当前文件中的类
            if seg != class_name:
                if seg in self.gctx.classes:
                    ncursor = self.gctx.classes[seg]
                    seg_code_file = self.gutil.get_code_file(ncursor)
                    cur_code_file = self.gutil.get_code_file(cursor)
                    seg_mod = self.gutil.get_decl_mod(ncursor)
                    cur_mod = self.gutil.get_decl_mod(cursor)

                    if seg_mod != cur_mod:  # 引用的类不在当前mod中
                        ctx.CP.APU('use', "use super::super::%s::%s::%s; // 771" % (seg_mod, seg_code_file, seg))
                    else:
                        if seg_code_file == cur_code_file:
                            ctx.CP.APU('use', "// use super::%s::%s; // 773" % (seg_code_file, seg))
                        else:
                            ctx.CP.APU('use', "use super::%s::%s; // 773" % (seg_code_file, seg))
                else:
                    # 不在类列表中的引用不了，如果有使用的地方，还是再找原因比较好
                    ctx.CP.APU('use', "// use super::%s::%s; // 775" % (seg.lower(), seg))
        return

    def dedup_return_const_diff_method(self, methods):
        dupremove = []
        for mtop in methods:
            postop = mtop.find('Q')
            for msub in methods:
                if mtop == msub: continue
                possub = msub.find('Q')
                if mtop[postop:] != msub[possub:]: continue
                if postop > possub: dupremove.append(mtop)
                else: dupremove.append(msub)
        return dupremove

    def reform_return_type_name(self, retname):
        lst = retname.split(' ')
        for elem in lst:
            if self.is_qt_class(elem): return elem
            if elem == 'String': return elem
        return retname

    def is_qt_class(self, type_name):
        # should be qt class name
        for seg in type_name.split(' '):
            if seg[0:1] == 'Q' and seg[1:2].upper() == seg[1:2] and '::' not in seg:  # should be qt class name
                return True
        return False

    def get_qt_class(self, type_name):
        # should be qt class name
        for seg in type_name.split(' '):
            if seg[0:1] == 'Q' and seg[1:2].upper() == seg[1:2] and '::' not in seg:  # should be qt class name
                return seg
        return None

    def fix_conflict_method_name(self, method_name):
        mthname = method_name
        fixmthname = mthname
        if mthname in ['match', 'type', 'move']:  # , 'select']:
            fixmthname = mthname + '_'
        return fixmthname

    def is_conflict_method_name(self, method_name):
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

            # Fix && move语义参数方法，
            if '&&' in type_name: return True
            if arg.type.kind == clang.cindex.TypeKind.RVALUEREFERENCE: return True
            if 'QPrivate' in type_name: return True
            if 'Private' in type_name: return True
            if 'QAbstract' in type_name: return True
            if 'QLatin1String' == type_name: return True
            if 'QLatin1Char' == type_name: return True
            if 'QStringRef' in type_name: return True
            if 'QStringDataPtr' in type_name: return True
            if 'QByteArrayDataPtr' in type_name: return True
            if 'QModelIndexList' in type_name: return True
            if 'QXmlStreamNamespaceDeclarations' in type_name: return True
            if 'QGenericArgument' in type_name: return True
            if 'QJson' in type_name: return True
            # if 'QWidget' in type_name: return True
            if 'QTextEngine' in type_name: return True
            # if 'QAction' in type_name: return True
            if 'QPlatformPixmap' in type_name: return True
            if 'QPlatformScreen' in type_name: return True
            if 'QPlatformMenu' in type_name: return True
            if 'QFileDialogArgs' in type_name: return True
            if 'FILE' in type_name: return True
            if type_name[0:1] == 'Q' and '::' in type_name: return True  # 有可能是类内类，像QMetaObject::Connection
            if '<' in type_name: return True  # 模板类参数
            # void directoryChanged(const QString & path, QFileSystemWatcher::QPrivateSignal arg0);
            # 这个不准确，会把QCoreApplication(int &, char**, int)也过滤掉了
            if method_name == 'QCoreApplication':pass
            else:
                if arg.displayname == '' and type_name == 'int':
                    print(555, 'whyyyyyyyyyyyyyy', method_name, arg.type.spelling)
                    # return True  # 过滤的不对，前边的已经过滤掉。

            #### more
            can_type = self.tyconv.TypeToCanonical(arg.type)
            if can_type.kind == clang.cindex.TypeKind.FUNCTIONPROTO: return True
            # if method_name == 'fromRotationMatrix':
            if can_type.kind == clang.cindex.TypeKind.RECORD:
                decl = can_type.get_declaration()
                for token in decl.get_tokens():
                    # print(555, token.spelling)
                    if token.spelling == 'template': return True
                    break
                # print(555, can_type.kind, method_name, decl.kind, decl.spelling,
                      #decl.get_num_template_arguments(),
                      #)
                # exit(0)

        return False

    # @return True | False
    def check_skip_method(self, cursor):
        method_name = cursor.spelling
        if method_name.startswith('operator'):
            # print("Omited operator method: " + mth)
            return True

        # print('pub:' + str(cursor.access_specifier))
        if cursor.access_specifier == clang.cindex.AccessSpecifier.PUBLIC:
            pass
        if cursor.access_specifier == clang.cindex.AccessSpecifier.PROTECTED:
            if cursor.kind == clang.cindex.CursorKind.CONSTRUCTOR or \
               cursor.kind == clang.cindex.CursorKind.DESTRUCTOR:
                pass
            else: return True
        if cursor.access_specifier == clang.cindex.AccessSpecifier.PRIVATE:
            if cursor.kind == clang.cindex.CursorKind.CONSTRUCTOR or \
               cursor.kind == clang.cindex.CursorKind.DESTRUCTOR:
                pass
            else: return True

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

        #### toUpper() &&，c++11 move语义的方法去掉
        # _ZNKR7QString7toUpperEv, _ZNO7QString7toUpperEv
        mangled_name = cursor.mangled_name
        if mangled_name.startswith('_ZNO'): return True
        # TODO fix QString::data() vs. QString::data() const
        # _ZN7QString4dataEv, _ZNK7QString4dataEv
        # TODO 这种情况还挺多的。函数名相同，返回值不回的重载方法 。需要想办法处理。
        # 这是支持方式，http://stackoverflow.com/questions/24594374/overload-operators-with-different-rhs-type
        # widgets
        # if mangled_name == '_ZN8QMenuBar7addMenuEP5QMenu': return True
        # if mangled_name == '_ZN5QMenu7addMenuEPS_': return True
        # if mangled_name == '_ZN11QMainWindow10addToolBarEP8QToolBar': return True
        # if mangled_name == '_ZNK15QCalendarWidget14dateTextFormatEv': return True
        # if mangled_name == '_ZN9QScroller8scrollerEPK7QObject': return True
        # if mangled_name == '_ZN12QApplication8setStyleERK7QString': return True
        # if method_name == 'mapToScene': return True  # 重载的方法太多
        # if method_name == 'mapFromScene': return True  # 重载的方法太多
        # if method_name == 'mapToItem': return True
        # if method_name == 'mapToParent': return True
        # if method_name == 'mapFromItem': return True
        # if method_name == 'mapFromParent': return True
        # if method_name == 'resolve': return True  # QFont::resolve
        # if method_name == 'map': return True  # QTransform::map
        # if method_name == 'mapRect': return True  # QTransform::mapRect
        # if method_name == 'point': return True  # QPolygon::point
        # if method_name == 'boundingRect': return True  # QPainter::boundingRect
        # if method_name == 'borderColor': return True  # QOpenGLTexture::borderColor
        # if method_name == 'trueMatrix': return True  # QPixmap::trueMatrix
        # if method_name == 'insertRow': return True  # QStandardItemModel::insertRow
        # gui
        class_name = cursor.semantic_parent.spelling
        # if method_name == 'read' and class_name == 'QImageReader': return True
        # if method_name == 'find' and class_name == 'QPixmapCache': return True
        # core
        # if class_name == 'QChar' and method_name == 'toUpper': return True
        # if class_name == 'QChar' and method_name == 'toLower': return True
        # if class_name == 'QChar' and method_name == 'mirroredChar': return True
        # if class_name == 'QChar' and method_name == 'toTitleCase': return True
        # if class_name == 'QChar' and method_name == 'toCaseFolded': return True
        # if class_name == 'QByteArray' and method_name == 'fill': return True
        # if class_name == 'QBitArray' and method_name == 'fill': return True
        # if class_name == 'QIODevice' and method_name == 'read': return True
        # if class_name == 'QIODevice' and method_name == 'peek': return True
        # if class_name == 'QIODevice' and method_name == 'readLine': return True
        # if class_name == 'QFileSelector' and method_name == 'select': return True
        # if class_name == 'QTextDecoder' and method_name == 'toUnicode': return True
        # if class_name == 'QCryptographicHash' and method_name == 'addData': return True
        # if class_name == 'QMessageAuthenticationCode' and method_name == 'addData': return True

        # 实现不知道怎么fix了，已经fix，原来是给clang.cindex.parse中的-I不全，导致找不到类型。
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

        # like 可能是Qt类内enum
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

        if atype.kind == clang.cindex.TypeKind.TYPEDEF:
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
            if type_bclass.kind == clang.cindex.CursorKind.TRANSLATION_UNIT: pass
            else: type_name = '%s::%s' % (type_bclass.spelling, type_name)
        else:
            type_name = self.real_type_name(atype)

            # QTextStreamManipulator(void (QTextStream::*)(int) m, int a);
            # int registerNormalizedType(const ::QByteArray & normalizedTypeName, void * destructor, void *(*)(void *, const void *) constructor, int size, QMetaType::TypeFlags flags, const QMetaObject * metaObject);
            # qreal (*)(qreal) customType();
            # if type_name == 'void (*)(void *)':
            #    type_name = "void *"

        return type_name

    def get_cursor_tokens(self, cursor):
        tokens = []
        for token in cursor.get_tokens():
            tokens.append(token.spelling)
        return ' '.join(tokens)

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
            print('write mod.rs:', mod, len(code), lines)
            self.write_modrs(mod, code)
        return

    def write_code(self, mod, fname, code):
        # mod = 'core'
        # fpath = "src/core/%s.rs" % (fname)
        fpath = "src/%s/%s.rs" % (mod, fname)
        self.write_file(fpath, code)
        return

    # TODO dir is exists
    def write_file(self, fpath, code):
        f = os.open(fpath, os.O_CREAT | os.O_TRUNC | os.O_RDWR)
        os.write(f, code)
        os.close(f)

        return

    def write_modrs(self, mod, code):
        fpath = "src/%s/mod.rs" % (mod)
        self.write_file(fpath, code)
        return
    pass


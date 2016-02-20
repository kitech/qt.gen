# encoding: utf8

import os
import logging
import time

import clang
import clang.cindex
import clang.cindex as clidx

from genutil import *
from typeconv import TypeConv, TypeConvForRust
from typeconvgo import TypeConvForGo
from genbase import GenerateBase, GenClassContext, GenMethodContext
from genfilter import GenFilterGo


class GenerateForGo(GenerateBase):
    def __init__(self):
        super(GenerateForGo, self).__init__()

        self.gfilter = GenFilterGo()
        self.modrss = {}  # mod => CodePaper
        #self.cp_modrs = CodePaper()  # 可能的name: main
        #self.cp_modrs.addPoint('main')
        #self.MP = self.cp_modrs

        self.class_blocks = ['header', 'main', 'use', 'ext', 'body']
        # self.cp_clsrs = CodePaper()  # 可能中间reset。可能的name: header, main, use, ext, body
        # self.CP = self.cp_clsrs

        self.qclses = {}  # class name => True
        self.tyconv = TypeConvForRust()
        self.tyconv = TypeConvForGo()
        self.traits = {}  # traits proto => True
        self.implmthods = {}  # method proto => True
        return

    def generateHeader(self, module):
        code = ''
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
            # CP.AP('header', 'package %s' % ('qt5'))
            CP.AP('header', 'package qt%s' % (mod))
            CP.AP('header', '// auto generated, do not modify.')
            CP.AP('header', '// created: ' + time.ctime())
            CP.AP('header', '// src-file: ' + key)
            CP.AP('header', '// dst-file: /src/%s/%s.go' % (mod, code_file))
            CP.AP('header', '//\n')

            for blk in self.class_blocks:
                CP.addPoint(blk)
                CP.append(blk, "// %s block begin =>" % (blk))
        return

    def genpass_code_header(self):
        modeps = {'core': [], 'gui': ['core'], 'widgets': ['core', 'gui']}

        for key in self.gctx.codes:
            CP = self.gctx.codes[key]
            CP.AP('header', self.generateHeader(''))
            CP.AP('use', 'import "fmt"')
            CP.AP('use', 'import "reflect"')
            CP.AP('use', 'import "unsafe"')
            CP.AP('use', 'import "qtrt"')
            mod = self.gctx.get_decl_mod_by_path(key)
            for dep in modeps[mod]:
                CP.AP('use', 'import "qt%s"' % (dep))

            # CP.AP('ext', "// extern {")
            CP.AP('ext', "\n/*")
            CP.AP('ext', "#include <stdlib.h>")
            CP.AP('ext', "#include <stdbool.h>")
            CP.AP('ext', "#include <stdint.h>")
            CP.AP('ext', "#include <wchar.h>")
            CP.AP('ext', "#include <uchar.h>")

            CP.AP('body', 'func init() {')
            CP.AP('body', '  if false {qtrt.KeepMe()}')
            for dep in modeps[mod]:
                CP.AP('body', '  if false {qt%s.KeepMe()}' % (dep))
            CP.AP('body', '  if false {fmt.Println(123)}')
            CP.AP('body', '  if false {reflect.TypeOf(123)}')
            CP.AP('body', '  if false {reflect.TypeOf(unsafe.Sizeof(0))}')
            CP.AP('body', '}\n')
        return

    def genpass_code_endian(self):
        for key in self.gctx.codes:
            CP = self.gctx.codes[key]
            for blk in self.class_blocks:
                if blk == 'ext':
                    CP.append(blk, '*/\nimport "C"')
                    CP.append(blk, "// } // <= %s block end\n" % (blk))
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

        print('gen signals...')
        # self.genpass_classes_signals()

        print('gen code endian...')
        self.genpass_code_endian()

        print('gen class mod define...')
        # self.genpass_class_modef()

        print('gen files...')
        self.genpass_write_codes()
        return

    def genpass_class_type(self):
        for key in self.gctx.classes:
            cursor = self.gctx.classes[key]
            if self.check_skip_class(cursor): continue
            self.genpass_class_type_impl(cursor)
        return

    def flat_template_name(self, name):
        return self.gutil.flat_template_name(name)

    def genpass_class_type_impl(self, cursor):
        class_name = cursor.type.spelling
        flat_class_name = self.flat_template_name(class_name)
        decl_file = self.gctx.get_decl_file(cursor)
        CP = self.gctx.codes[decl_file]
        ctysz = cursor.type.get_size()

        # TODO 计算了两遍
        bases = self.gutil.get_base_class(cursor)
        base_class = bases[0] if len(bases) > 0 else None
        usignals = self.gutil.get_unique_signals(cursor)

        CP.AP('body', "// class sizeof(%s)=%s" % (class_name, ctysz))
        # generate struct of class
        # CP.AP('body', '#[derive(Sized)]')
        # CP.AP('body', '#[derive(Default)]')
        CP.AP('body', "type %s struct {" % (flat_class_name))
        if base_class is None:
            CP.AP('body', "  // qbase: %s;" % (base_class))
        else:
            # TODO 需要use 基类
            bmod = self.gutil.get_decl_mod(self.get_qt_class_cursor(base_class.spelling))
            cmod = self.gutil.get_decl_mod(cursor)
            if bmod != cmod:
                CP.AP('body', "  /*qbase*/ qt%s.%s;" % (bmod, base_class.spelling))
            else:
                CP.AP('body', "  /*qbase*/ %s;" % (base_class.spelling))
        CP.AP('body', "  Qclsinst unsafe.Pointer /* *C.void */;")
        for key in usignals:
            sigmth = usignals[key]
            CP.AP('body', '//  _%s %s_%s_signal;' % (sigmth.spelling, flat_class_name, sigmth.spelling))
        CP.AP('body', "}\n")

        return

    def genpass_classes_signals(self):
        for key in self.gctx.classes:
            cursor = self.gctx.classes[key]
            if self.check_skip_class(cursor): continue

            class_name = cursor.displayname
            methods = self.gutil.get_methods(cursor)
            bases = self.gutil.get_base_class(cursor)
            base_class = bases[0] if len(bases) > 0 else None
            ctx = mctx = self.createMiniContext(cursor, base_class)
            usignals = self.gutil.get_unique_signals(cursor)
            for key in usignals:
                sigmth = usignals[key]
                ctx.CP.AP('body', '#[derive(Default)] // for %s_%s' % (class_name, sigmth.spelling))
                ctx.CP.AP('body', 'pub struct %s_%s_signal{poi:u64}' % (class_name, sigmth.spelling))

                ctx.CP.AP('body', 'impl /* struct */ %s {' % (class_name))
                ctx.CP.AP('body', '  pub fn %s(&self) -> %s_%s_signal {'
                          % (sigmth.spelling, class_name, sigmth.spelling))
                ctx.CP.AP('body', '     return %s_%s_signal{poi:self.Qclsinst};' % (class_name, sigmth.spelling))
                ctx.CP.AP('body', '  }')
                ctx.CP.AP('body', '}')

                ctx.CP.AP('body', 'impl /* struct */ %s_%s_signal {' % (class_name, sigmth.spelling))
                ctx.CP.AP('body', '  pub fn connect<T: %s_%s_signal_connect>(self, overload_args: T) {'
                          % (class_name, sigmth.spelling))
                ctx.CP.AP('body', '    overload_args.connect(self);')
                ctx.CP.AP('body', '  }')
                ctx.CP.AP('body', '}')
                ctx.CP.AP('body', 'pub trait %s_%s_signal_connect {' % (class_name, sigmth.spelling))
                ctx.CP.AP('body', '  fn connect(self, sigthis: %s_%s_signal);' % (class_name, sigmth.spelling))
                ctx.CP.AP('body', '}')
                ctx.CP.AP('body', '')

            idx = 0
            for key in ctx.signals:
                sigmth = ctx.signals[key]
                if '<' in sigmth.displayname: continue
                ctx = self.createGenMethodContext(sigmth, cursor, base_class, [])

                trait_params_array = self.generateParamsForTrait(class_name, sigmth.spelling, sigmth, ctx)
                trait_params = ', '.join(trait_params_array)
                trait_params = trait_params.replace("&'a mut ", '')
                trait_params = trait_params.replace("&'a ", '')
                if '<' in trait_params: continue  # QModelIndexList => QList<QModelIndex>
                if 'QPrivateSignal' in trait_params: continue

                params_ext_arr = self.generateParamsForExtern(class_name, sigmth.spelling, sigmth, ctx)
                params_ext = ', '.join(params_ext_arr)
                params_ext_tyarr = []
                for arg in params_ext_arr:
                    params_ext_tyarr.append(arg.split(':')[1].strip())
                params_ext_ty = ', ' .join(params_ext_tyarr)

                ctx.CP.AP('body', '// %s' % (sigmth.displayname))
                ctx.CP.AP('body', 'extern fn %s_%s_signal_connect_cb_%s(rsfptr:fn(%s), %s) {'
                          % (class_name, sigmth.spelling, idx, trait_params, params_ext))
                ctx.CP.AP('body', '  println!("{}:{}", file!(), line!());')
                rsargs = []
                for arg in sigmth.get_arguments():
                    sidx = len(rsargs)
                    sty = params_ext_arr[sidx]
                    dty = trait_params_array[sidx]
                    if self.is_qt_class(arg.type.spelling):
                        arg_class_name = self.get_qt_class(arg.type.spelling)
                        ctx.CP.AP('body', '  let rsarg%s = %s::inheritFrom(arg%s as u64);'
                                  % (sidx, arg_class_name, sidx))
                    else:
                        ctx.CP.AP('body', '  let rsarg%s = arg%s as %s;' % (sidx, sidx, dty))
                    rsargs.append('rsarg%s' % (sidx))

                ctx.CP.AP('body', '  rsfptr(%s);' % (','.join(rsargs)))
                ctx.CP.AP('body', '}')
                ctx.CP.AP('body', 'extern fn %s_%s_signal_connect_cb_box_%s(rsfptr_raw:*mut Box<Fn(%s)>, %s) {'
                          % (class_name, sigmth.spelling, idx, trait_params, params_ext))
                ctx.CP.AP('body', '  println!("{}:{}", file!(), line!());')
                ctx.CP.AP('body', '  let rsfptr = unsafe{Box::from_raw(rsfptr_raw)};')

                rsargs = []
                for arg in sigmth.get_arguments():
                    sidx = len(rsargs)
                    sty = params_ext_arr[sidx]
                    dty = trait_params_array[sidx]
                    if self.is_qt_class(arg.type.spelling):
                        arg_class_name = self.get_qt_class(arg.type.spelling)
                        ctx.CP.AP('body', '  let rsarg%s = %s::inheritFrom(arg%s as u64);'
                                  % (sidx, arg_class_name, sidx))
                    else:
                        ctx.CP.AP('body', '  let rsarg%s = arg%s as %s;' % (sidx, sidx, dty))
                    rsargs.append('rsarg%s' % (sidx))

                ctx.CP.AP('body', '  // rsfptr(%s);' % (','.join(rsargs)))
                ctx.CP.AP('body', '  unsafe{(*rsfptr_raw)(%s)};' % (','.join(rsargs)))
                ctx.CP.AP('body', '}')
                # impl xxx for fn(%s)
                ctx.CP.AP('body', 'impl /* trait */ %s_%s_signal_connect for fn(%s) {'
                          % (class_name, sigmth.spelling, trait_params))
                ctx.CP.AP('body', '  fn connect(self, sigthis: %s_%s_signal) {' % (class_name, sigmth.spelling))
                ctx.CP.AP('body', '    // do smth...')
                ctx.CP.AP('body', '    // self as u64; // error for Fn, Ok for fn')
                ctx.CP.AP('body', '    self as *mut c_void as u64;')
                ctx.CP.AP('body', '    self as *mut c_void;')
                ctx.CP.AP('body', '    let arg0 = sigthis.poi as *mut c_void;')
                ctx.CP.AP('body', '    let arg1 = %s_%s_signal_connect_cb_%s as *mut c_void;'
                          % (class_name, sigmth.spelling, idx))
                ctx.CP.AP('body', '    let arg2 = self as *mut c_void;')
                # ctx.CP.AP('body', '    // %s_%s_signal_connect_cb_%s' % (class_name, sigmth.spelling, idx))
                ctx.CP.AP('body', '    unsafe {%s_SlotProxy_connect_%s(arg0, arg1, arg2)};'
                          % (class_name, sigmth.mangled_name))
                ctx.CP.AP('body', '  }')
                ctx.CP.AP('body', '}')
                # impl xx for Box<fn(%s)>
                ctx.CP.AP('body', 'impl /* trait */ %s_%s_signal_connect for Box<Fn(%s)> {'
                          % (class_name, sigmth.spelling, trait_params))
                ctx.CP.AP('body', '  fn connect(self, sigthis: %s_%s_signal) {' % (class_name, sigmth.spelling))
                ctx.CP.AP('body', '    // do smth...')
                ctx.CP.AP('body', '    // Box::into_raw(self) as u64;')
                ctx.CP.AP('body', '    // Box::into_raw(self) as *mut c_void;')
                ctx.CP.AP('body', '    let arg0 = sigthis.poi as *mut c_void;')
                ctx.CP.AP('body', '    let arg1 = %s_%s_signal_connect_cb_box_%s as *mut c_void;'
                          % (class_name, sigmth.spelling, idx))
                ctx.CP.AP('body', '    let arg2 = Box::into_raw(Box::new(self)) as *mut c_void;')
                # ctx.CP.AP('body', '    // %s_%s_signal_connect_cb_%s' % (class_name, sigmth.spelling, idx))
                ctx.CP.AP('body', '    unsafe {%s_SlotProxy_connect_%s(arg0, arg1, arg2)};'
                          % (class_name, sigmth.mangled_name))
                ctx.CP.AP('body', '  }')
                ctx.CP.AP('body', '}')
                ctx.CP.AP('ext', '  fn %s_SlotProxy_connect_%s(qthis: *mut c_void, ffifptr: *mut c_void, rsfptr: *mut c_void);'
                          % (class_name, sigmth.mangled_name))
                idx += 1

            # self.generateClass(class_name, cursor, methods, base_class)
            # break
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

        signals = self.gutil.get_signals(class_cursor)
        class_overload_methods = self.overload_method_groupby(methods)
        for method_name in class_overload_methods:
            mangled_names = self.overload_method_select(class_overload_methods[method_name])
            method_overload_methods = {}  # mangled_name => cursor
            for mangled_name in mangled_names:
                cursor = methods[mangled_name]
                if self.check_skip_method(cursor): continue
                if self.check_skip_params(cursor): continue

                if mangled_name in signals: continue
                method_overload_methods[mangled_name] = cursor

            if len(method_overload_methods) == 0: continue
            self.generateMethod(method_name, method_overload_methods, class_cursor, base_class, unique_methods)

        return

    def createGenMethodContext(self, method_cursor, class_cursor, base_class, unique_methods):
        ctx = GenMethodContext(method_cursor, class_cursor)
        ctx.unique_methods = unique_methods
        ctx.CP = self.gctx.getCodePager(class_cursor)

        if ctx.ctor: ctx.method_name_rewrite = 'New%s' % (ctx.class_name)
        if ctx.dtor: ctx.method_name_rewrite = 'Free%s' % (ctx.class_name)
        if self.is_conflict_method_name(ctx.method_name):
            ctx.method_name_rewrite = ctx.method_name + '_'
        if ctx.static:
            ctx.method_name_rewrite = ctx.method_name + ctx.static_suffix

        ctx.isinline = self.method_is_inline(method_cursor)

        class_name = ctx.class_name
        method_name = ctx.method_name


        # ctx.ret_type_name_rs = self.tyconv.Type2RustRet(ctx.ret_type, method_cursor)
        ctx.ret_type_name_ext = self.tyconv.ArgType2FFIExt(ctx.ret_type, method_cursor)

        raw_params_array = self.generateParamsRaw(class_name, method_name, method_cursor)
        raw_params = ', '.join(raw_params_array)

        # trait_params_array = self.generateParamsForTrait(class_name, method_name, method_cursor, ctx)
        # trait_params = ', '.join(trait_params_array)

        call_params_array = self.generateParamsForCall(class_name, method_name, method_cursor)
        # if ctx.ctor: call_params_array.insert(0, 'qthis')
        call_params = ', '.join(call_params_array)
        if not ctx.static and not ctx.ctor: call_params = ('this.Qclsinst, ' + call_params).strip(' ,')

        extargs_array = self.generateParamsForExtern(class_name, method_name, method_cursor, ctx)
        extargs = ', '.join(extargs_array)
        if not ctx.static and not ctx.ctor: extargs = ('void* qthis, ' + extargs).strip(' ,')

        ctx.params_cpp = raw_params
        # ctx.params_rs = trait_params
        ctx.params_call = call_params
        ctx.params_ext = extargs
        ctx.params_ext_arr = extargs_array

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
        # ctx.tymap = TypeConvForRust.tymap

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

    def generateMethod(self, method_name, overload_methods, class_cursor, base_classes, unique_methods):
        first_method_cursor = overload_methods.values()[0]
        ctx = self.createGenMethodContext(first_method_cursor, class_cursor, base_classes, unique_methods)
        flat_class_name = self.flat_template_name(class_cursor.type.spelling)
        flat_method_name = self.flat_template_name(ctx.method_name_rewrite)

        ctx.CP.AP('body', '// %s' % (first_method_cursor.displayname))
        if ctx.ctor:
            ctx.CP.AP('body', 'func %s(args ...interface{}) *%s {'
                      % (flat_method_name, flat_class_name))
        else:
            if ctx.has_return:
                ctx.CP.AP('body', 'func (this *%s) %s(args ...interface{}) (ret interface{}) {'
                          % (flat_class_name, flat_method_name.title()))

            else:
                ctx.CP.AP('body', 'func (this *%s) %s(args ...interface{}) () {'
                          % (flat_class_name, flat_method_name.title()))

        for mangled_name in overload_methods:
            cursor = overload_methods[mangled_name]
            ctx.CP.AP('body', "  // %s" % (cursor.displayname))

        ctx.CP.AP('body', '  var vtys = make(map[int32]map[int32]reflect.Type)')
        ctx.CP.AP('body', "  if false {fmt.Println(vtys)}")

        midx = 0
        for mangled_name in overload_methods:
            cursor = overload_methods[mangled_name]
            ctx.CP.AP('body', '  vtys[%s] = make(map[int32]reflect.Type)' % (midx))
            self.generateParamsTypeForResolve(ctx, cursor, midx)
            midx += 1

        ctx.CP.AP('body', '')
        ctx.CP.AP('body', "  var matched_index = qtrt.SymbolResolve(args, vtys)")
        ctx.CP.AP('body', "  if false {fmt.Println(matched_index)}")

        self.generateVTableInvoke(ctx, overload_methods)

        if ctx.ctor:
            # ctx.CP.AP('body', '  return %s{Qclsinst:qthis}' % (flat_class_name))
            ctx.CP.AP('body', '  return nil // %s{Qclsinst:qthis}' % (flat_class_name))
        else:
            ctx.CP.AP('body', '  return')
        ctx.CP.AP('body', "}\n")

        # extern
        for mangled_name in overload_methods:
            cursor = overload_methods[mangled_name]
            ctx = self.createGenMethodContext(cursor, class_cursor, base_classes, unique_methods)
            ctx.CP.AP('ext', ctx.fn_proto_cpp)
            self.generateDeclForFFIExt(ctx)

        return

    def generateVTableInvoke(self, ctx, overload_methods):
        movs = {}
        for mangled_name in overload_methods:
            movs[mangled_name] = overload_methods[mangled_name]
        # dedup_methods = self.dedup_return_const_diff_method(movs)
        midx = -1
        ctx.CP.AP('body', '  switch matched_index {')
        for mangled_name in overload_methods:
            midx += 1
            cursor = mth = overload_methods[mangled_name]
            # deprefix = 'dector' if ctx.ctor else 'demth'
            # deprefix = ''
            ctx.CP.AP('body', '  case %s:' % (midx))
            ctx.CP.AP('body', '    // invoke: %s' % (mth.mangled_name))
            ctx.CP.AP('body', '    // invoke: %s %s' % (mth.result_type.spelling, mth.displayname))
            nctx = self.createGenMethodContext(mth, ctx.class_cursor, ctx.base_class, ctx.unique_methods)
            self.generateArgConvExprs(ctx.class_name, mth.spelling, mth, nctx, midx)
            if nctx.ctor:
                ctx.CP.AP('body', '    var qthis = unsafe.Pointer(C.malloc(5))')
                ctx.CP.AP('body', '    if false {reflect.TypeOf(qthis)}')
                ctx.CP.AP('body', '    qthis = C.%s(%s)' % (nctx.cmangled_name, nctx.params_call))
                ctx.CP.AP('body', '    return &%s{Qclsinst:qthis}' % (nctx.flat_class_name))
            else:
                if nctx.has_return:
                    ctx.CP.AP('body', '    var ret0 = C.%s(%s)' % (nctx.cmangled_name, nctx.params_call))
                    ctx.CP.AP('body', '    if false {reflect.TypeOf(ret0)}')
                    # 竟然还有重载的 方c法，有的有返回值，有的没有
                    if ctx.has_return:
                        ctx.CP.AP('body', '    ret = ret0')
                        self.generateReturnTypeDecl(ctx, nctx)
                        self.generateReturn(ctx, nctx)
                else:
                    ctx.CP.AP('body', '    C.%s(%s)' % (nctx.cmangled_name, nctx.params_call))
            # if nctx.isinline:
            #     ctx.CP.AP('body', '    C.%s%s(%s)' % (deprefix, nctx.cmangled_name, nctx.params_call))
            # else:
            #     ctx.CP.AP('body', '    C.%s(%s)' % (nctx.cmangled_name, nctx.params_call))

        ctx.CP.AP('body', '  default:')
        ctx.CP.AP('body', '    qtrt.ErrorResolve("%s", "%s", args)' % (ctx.class_name, ctx.method_name))
        ctx.CP.AP('body', '  }\n')
        return

    def generateArgConvExprs(self, class_name, method_name, method_cursor, ctx, midx):
        argc = 0
        for arg in method_cursor.get_arguments(): argc += 1

        def isvec(tyname): return 'Vec<' in tyname
        def isrstr(tyname): return 'String' in tyname.split(' ')

        for idx, (arg) in enumerate(method_cursor.get_arguments()):
            atc = self.tyconv.ArgType2CGO(arg.type, arg)
            if '%s' not in atc:
                print(123, atc)
                raise '123'
            if atc.startswith('qtrt.HandyConvert2c'):
                atc = atc % ('args[%s]' % (idx), 'vtys[%s][%s]' % (midx, idx))
                if '<' in atc or '::' in atc:
                    ctx.CP.AP('body', "    // var arg%s = %s" % (idx, atc))
                    ctx.CP.AP('body', "    var arg%s unsafe.Pointer" % (idx))
                    ctx.CP.AP('body', '    if false {fmt.Println(arg%s)}' % (idx))
                else:
                    ctx.CP.AP('body', "    argif%s, free%s := %s" % (idx, idx, atc))
                    ctx.CP.AP('body', "    var arg%s = argif%s.(unsafe.Pointer)" % (idx, idx))
                    ctx.CP.AP('body', '    if false {fmt.Println(argif%s, arg%s)}' % (idx, idx))
                    ctx.CP.AP('body', '    if free%s {defer C.free(arg%s)}' % (idx, idx))
            else:
                atc = atc % ('args[%s]' % (idx))
                if '<' in atc or '::' in atc:
                    ctx.CP.AP('body', "    // var arg%s = %s" % (idx, atc))
                    ctx.CP.AP('body', "    var arg%s unsafe.Pointer" % (idx))
                    ctx.CP.AP('body', '    if false {fmt.Println(arg%s)}' % (idx))
                else:
                    ctx.CP.AP('body', "    var arg%s = %s" % (idx, atc))
                    ctx.CP.AP('body', '    if false {fmt.Println(arg%s)}' % (idx))
        return

    def generateParamsTypeForResolve(self, ctx, method_cursor, method_index):
        midx = method_index
        for idx, (arg) in enumerate(method_cursor.get_arguments()):
            arty = self.tyconv.ArgType2GoReflectType(arg.type, arg)
            # print(345, arg.type.spelling, '=>', arty)
            if '<' in arty or '::' in arty:
                ctx.CP.AP('body', '  // vtys[%s][%s] = %s // "%s"' % (midx, idx, arty, arg.type.spelling))
            else:
                ctx.CP.AP('body', '  vtys[%s][%s] = %s // "%s"' % (midx, idx, arty, arg.type.spelling))
        return

    # @return []
    def generateParamsRaw(self, class_name, method_name, method_cursor):
        argv = []
        for arg in method_cursor.get_arguments():
            argelem = "%s %s" % (arg.type.spelling, arg.displayname)
            argv.append(argelem)
        return argv

    # @return []
    def generateParamsForCall(self, class_name, method_name, method_cursor):
        idx = -1
        argv = []

        for arg in method_cursor.get_arguments():
            idx += 1
            # print('%s, %s, ty:%s, kindty:%s' % (method_name, arg.displayname, arg.type.spelling, arg.kind))
            # print('arg type kind: %s, %s' % (arg.type.kind, arg.type.get_declaration()))
            argelem = 'arg%s' % (idx)
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

            if self.check_skip_param(arg, method_name) is False:
                self.generateUseForRust(ctx, arg.type, arg)

            type_name = self.tyconv.TypeCXX2Rust(arg.type, arg, inty=True)
            if type_name.startswith('&'): type_name = type_name.replace('&', "&'a ")

            arg_name = 'arg%s' % idx if arg.displayname == '' else arg.displayname
            argelem = "%s" % (type_name)
            argv.append(argelem)

        return argv

    # @return []
    def generateParamsForExtern(self, class_name, method_name, method_cursor, ctx):
        idx = 0
        argv = []

        for arg in method_cursor.get_arguments():
            idx += 1
            # print('%s, %s, ty:%s, kindty:%s' % (method_name, arg.displayname, arg.type.spelling, arg.kind))
            # print('arg type kind: %s, %s' % (arg.type.kind, arg.type.get_declaration()))
            type_name = self.tyconv.ArgType2FFIExt(arg.type, arg)
            arg_name = 'arg%s' % idx if arg.displayname == '' else arg.displayname
            argelem = "%s arg%s" % (type_name, idx - 1)
            argv.append(argelem)

        return argv

    def generateDeclForFFIExt(self, ctx):
        cursor = ctx.cursor
        has_return = ctx.has_return
        # calc ext type name
        return_type_name = 'void'
        # return_type_name = self.tyconv.TypeCXX2RustExtern(ctx.ret_type, cursor)
        return_type_name = ctx.ret_type_name_ext

        mangled_name = ctx.cmangled_name
        return_piece_proto = 'void'
        if cursor.result_type.kind != clidx.TypeKind.VOID and has_return:
            return_piece_proto = '%s' % (return_type_name)
        extargs = ctx.params_ext

        if ctx.isinline:
            if ctx.ctor:
                ctx.CP.AP('ext', "extern void* %s(%s); // 1" % (mangled_name, extargs))
            else:
                ctx.CP.AP('ext', "extern %s %s(%s); // 2" % (return_piece_proto, mangled_name, extargs))
        else:
            if ctx.ctor:
                ctx.CP.AP('ext', "extern void* %s(%s); // 3" % (mangled_name, extargs))
            else:
                ctx.CP.AP('ext', "extern %s %s(%s); // 4" % (return_piece_proto, mangled_name, extargs))

        return has_return, return_type_name

    def generateReturnTypeDecl(self, ctx, nctx):
        if ctx.has_return:
            arty = self.tyconv.ArgType2GoReflectType(nctx.cursor.result_type, nctx.cursor)
            if '<' in arty or '::' in arty:
                ctx.CP.AP('body', '    // var rety = %s // "%s"' % (arty, nctx.ret_type_name_cpp))
            else:
                ctx.CP.AP('body', '    var rety = %s // "%s"' % (arty, nctx.ret_type_name_cpp))
        return

    def generateReturn(self, ctx, nctx):
        if ctx.has_return:
            ctx.CP.AP('body', '    if reflect.TypeOf(ret0).ConvertibleTo(rety) {')
            ctx.CP.AP('body', '        ret = reflect.ValueOf(ret0).Convert(rety).Interface()')
            ctx.CP.AP('body', '    } else {')
            ctx.CP.AP('body', '        ret = qtrt.HandyConvert2go(ret0, rety)')
            ctx.CP.AP('body', '    }')
        return

    def methodHasReturn(self, ctx):
        method_cursor = cursor = ctx.cursor
        class_name = ctx.class_name

        return_type = cursor.result_type

        return_type_name = return_type.spelling
        if ctx.ctor or ctx.dtor: pass
        else:
            return_type_name = self.resolve_swig_type_name(class_name, return_type)

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

    def overload_method_groupby(self, methods):
        overloads = {}
        for mn in methods:
            name = methods[mn].spelling
            if name not in overloads:
                overloads[name] = []
            overloads[name].append(mn)
        return overloads

    # @param []
    def overload_method_select(self, overload_methods):
        if len(overload_methods) == 1: return overload_methods

        duprem = []
        # && (), () const, ()
        for mn in overload_methods:
            # drop && return if has dup
            if mn.startswith('_ZNO'):
                if mn.replace('_ZNO', '_ZNKR') in overload_methods:
                    duprem.append(mn)
                    continue
            # drop const return if has dup
            if mn.startswith('_ZNK'):
                if mn.replace('_ZNK', '_ZN') in overload_methods:
                    duprem.append(mn)
                    continue

        save = []
        for mn in overload_methods:
            if mn not in duprem:
                save.append(mn)

        # print(overload_methods, '=>', save)
        return save

    # TODO depcreated
    def is_qt_class(self, type_name):
        # should be qt class name
        for seg in type_name.split(' '):
            if seg[0:1] == 'Q' and seg[1:2].upper() == seg[1:2] and '::' not in seg:  # should be qt class name
                return True
        return False

    # TODO depcreated
    def get_qt_class(self, type_name):
        # should be qt class name
        for seg in type_name.split(' '):
            if seg[0:1] == 'Q' and seg[1:2].upper() == seg[1:2] and '::' not in seg:  # should be qt class name
                return seg
        return None

    def fix_conflict_method_name(self, method_name):
        mthname = method_name
        fixmthname = mthname
        if mthname in ['match', 'type', 'move', 'select', 'map']:
            fixmthname = mthname + '_'
        return fixmthname

    def is_conflict_method_name(self, method_name):
        if method_name in ['match', 'type', 'move', 'select', 'map']:
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
            # if 'QAbstract' in type_name: return True
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
            if type_name[0] == 'Q' and '::' in type_name: return True  # 有可能是类内类，像QMetaObject::Connection
            if '<' in type_name: return True  # 模板类参数
            if type_name[0] == 'Q' and type_name.endswith('Private'): return True

            # void directoryChanged(const QString & path, QFileSystemWatcher::QPrivateSignal arg0);
            # 这个不准确，会把QCoreApplication(int &, char**, int)也过滤掉了
            if method_name == 'QCoreApplication':pass
            else:
                if arg.displayname == '' and type_name == 'int':
                    # print(555, 'whyyyyyyyyyyyyyy', method_name, arg.type.spelling)
                    # return True  # 过滤的不对，前边的已经过滤掉。
                    pass

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
        # shitfix begin
        method_name = cursor.spelling
        mangled_name = cursor.mangled_name

        if 'QOpenGLFunctions_' in method_name: return True
        if 'QOpenGLFunctionsPrivate' == method_name: return True
        if 'QAbstractOpenGLFunctionsPrivate' == method_name: return True
        # if 'QTextStreamManipulator' == method_name: return True
        if 'QRunnable' == method_name: return True

        fixmths = ['data_ptr',
                   'sprintf', 'vsprintf', 'vasprintf', 'asprintf',
                   ]
        if method_name in fixmths: return True

        if method_name[0] == '~' and method_name.endswith('Interface'): return True
        if method_name[0] == 'Q' and method_name.endswith('Private'): return True

        howfixs = ['_ZN7QWindow9fromWinIdEi', '_ZN7QWidget4findEi',
                   '_ZN8QVariantC1EOS_', '_ZN4QUrlC1EOS_',
                   '_ZN14QTextTableCellD0Ev', '_ZN11QStringListC1EO5QListI7QStringE',
                   '_ZN7QString16fromStdU32StringERKi',
                   '_ZN7QString16fromStdU16StringERKi',
                   '_ZN7QString14fromStdWStringERKi', '_ZN7QString13fromStdStringERKi',
                   '_ZN15QSocketNotifierC1EiNS_4TypeEP7QObject', '_ZN9QRunnableD0Ev',
                   '_ZN7QPixmap9fromImageEO6QImage6QFlagsIN2Qt19ImageConversionFlagEE',
                   '_ZN7QPixmap10grabWindowEiiiii', '_ZN4QPenC1EOS_', '_ZN8QPaletteC1EOS_',
                   '_ZN11QMetaObject10ConnectionC1EOS0_', '_ZN14QSignalBlockerC1EOS_',
                   '_ZN6QImageC1EOS_', '_ZN12QEasingCurveC1EOS_', '_ZN7QCursorC1EOS_',
                   '_ZN9QCollatorC1EOS_', '_ZN10QByteArray13fromStdStringERKi',
                   '_ZN10QByteArrayC1EOS_', '_ZN9QBitArrayC1EOS_',
                   '_ZNK10QArrayData14detachCapacityEi',
                   '_ZN10QArrayData8allocateEiii6QFlagsINS_16AllocationOptionEE',
                   '_ZN10QArrayData10deallocateEPS_ii', '_ZN17QAccessibleBridgeD0Ev',
                   '_ZN20QAccessibleInterface14interface_castEN11QAccessible13InterfaceTypeE',
                   '_ZN21QPersistentModelIndexC1EOS_',
                   # link errors
                   '_ZN5QFile4openEP8_IO_FILE6QFlagsIN9QIODevice12OpenModeFlagEES2_IN11QFileDevice14FileHandleFlagEE',
                   '_ZN6QImageC1EPKPKc', '_ZN7QPixmapC1EPKPKc',
                   '_ZN5QMenu15setPlatformMenuEP13QPlatformMenu',
                   '_ZN7QPixmapC1EP15QPlatformPixmap',
                   '_ZN7QString9fromUtf16EPKDsi', '_ZN7QString8fromUcs4EPKDii',
                   '_ZN15QAnimationGroupC2EP7QObject', '_ZN17QAccessibleObjectC2EP7QObject',
                   ]
        if mangled_name in howfixs: return True

        absfixs = ['_ZN15QAnimationGroupC1EP7QObject', '_ZN17QAccessibleObjectC1EP7QObject',]
        if mangled_name in absfixs: return True

        # forward declaration type reference
        fdtfixs = ['_ZNK12QActionEvent6actionEv', '_ZNK12QActionEvent6beforeEv',
        ]
        if mangled_name in fdtfixs: return True

        # shitfix end

        if True: return self.gfilter.skipMethod(cursor)
        return False

    def check_skip_class(self, class_cursor):
        # shitfix begin

        cursor = class_cursor
        name = cursor.spelling
        dname = cursor.displayname

        if name.startswith('QOpenGLFunctions'): return True
        if name == 'QAbstractOpenGLFunctionsPrivate': return True
        if name == 'QSignalMapper': return True
        if name == 'QActionEvent': return True

        # shitfix end

        if True: return self.gfilter.skipClass(class_cursor)
        return False

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

        return type_name

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
            print('write mod.rs:', mod, len(code), lines)
            # self.write_modrs(mod, code)
        return

    def write_code(self, mod, fname, code):
        # mod = 'core'
        # fpath = "src/core/%s.rs" % (fname)
        fpath = "src/%s/%s.go" % (mod, fname)
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


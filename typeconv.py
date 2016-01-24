# encoding: utf8

import logging
import sys
import traceback

import clang.cindex as clidx

from genutil import *


# 类型转换规则
# 原始类型指针、引用
# 原始类型非指针、非引用
# char*类型
# qt类型指针，引用
# qt类型非指针，非引用

# 还需要分返回值的类型


# for python2 base class need to object
class TypeConvertContext(object):
    def __init__(self):
        self.orig_type = None
        self.can_type = None
        self.convable_type = None
        self.pointee_type = None
        self.pointer_level = 0
        self.const = False

        self.orig_type_name = ''
        self.can_type_name = ''
        self.convable_type_name = ''

        self.cursor = None
        return


class TypeConv(object):
    tymap = {}

    def __init__(self):
        return

    def TypeToCanonical(self, cxxtype):
        cxxtype = cxxtype.get_canonical()  # 这个一般不管用
        cxxtype.is_const_qualified()  # 这个函数也不管用

        if cxxtype.kind == clidx.TypeKind.POINTER or \
           cxxtype.kind == clidx.TypeKind.LVALUEREFERENCE or \
           cxxtype.kind == clidx.TypeKind.RVALUEREFERENCE:
            non_pointer_type = cxxtype.get_pointee()
            return self.TypeToCanonical(non_pointer_type)

        if cxxtype.kind == clidx.TypeKind.TYPEDEF:
            under_type = cxxtype.get_declaration().underlying_typedef_type
            return self.TypeToCanonical(under_type)

        return cxxtype

    def TypeToConvertable(self, cxxtype):
        if cxxtype.kind == clidx.TypeKind.TYPEDEF:
            under_type = cxxtype.get_declaration().underlying_typedef_type
            return self.TypeToConvertable(under_type)
        return cxxtype

    # 与TypeToCanonical不同，不解析指针
    # 不能失真的转换
    def TypeToActual(self, cxxtype):
        # cxxtype = cxxtype.get_canonical()  # 这个一般不管用
        cxxtype.is_const_qualified()  # 这个函数也不管用

        if cxxtype.kind == clidx.TypeKind.LVALUEREFERENCE:
            under_type = cxxtype.get_pointee()
            nty = self.TypeToActual(under_type)
            return nty

        if cxxtype.kind == clidx.TypeKind.RVALUEREFERENCE:
            under_type = cxxtype.get_pointee()
            nty = self.TypeToActual(under_type)
            return nty

        if cxxtype.kind == clidx.TypeKind.TYPEDEF:
            under_type = cxxtype.get_declaration().underlying_typedef_type
            nty = self.TypeToActual(under_type)
            return nty

        # 原来UNEXPOSED 就是编译时常见的未定义的 类型
        # TODO 查找一下哪些UNEXPOSED的
        if cxxtype.kind == clidx.TypeKind.UNEXPOSED:
            under_type = cxxtype.get_declaration().type
            nty = self.TypeToActual(under_type)
            return nty

        return cxxtype

    def TypeCanName(self, cxxtype):
        canty = self.TypeToCanonical(cxxtype)
        if canty.kind == clidx.TypeKind.FUNCTIONPROTO:
            # print(444, "function proto:", canty.spelling)
            # exit(0)
            pass
        canname = canty.spelling
        if self.TypeIsConst(canty): canname = self.TypeTrimConst(canty)
        if self.TypeIsVolatile(canty): canname = self.TypeTrimConst(canty)
        canname = canname.replace('unsigned ', 'u')
        return canname

    def TypeIsConst(self, cxxtype):
        return cxxtype.spelling.startswith('const ')

    def TypeIsVolatile(self, cxxtype):
        return cxxtype.spelling.startswith('volatile ')

    def TypeTrimConst(self, cxxtype):
        tysegs = cxxtype.spelling.split(' ')
        if 'const' in tysegs: tysegs.remove('const')
        if 'volatile' in tysegs: tysegs.remove('volatile')
        return ' '.join(tysegs)

    def TypeNameTrimConst(self, tyname):
        tysegs = tyname.split(' ')
        if 'const' in tysegs: tysegs.remove('const')
        if 'volatile' in tysegs: tysegs.remove('volatile')
        return ' '.join(tysegs)

    def IsCharType(self, tyname):
        if 'char' in tyname: return True
        if 'string' in tyname: return True  # for std::string, std::wstring
        return False

    def IsPointer(self, cxxtype):
        if cxxtype.kind == clidx.TypeKind.POINTER or \
           cxxtype.kind == clidx.TypeKind.LVALUEREFERENCE or \
           cxxtype.kind == clidx.TypeKind.RVALUEREFERENCE:
            return True
        return False

    def ArrayDim(self, cxxtype):
        name = cxxtype.spelling
        return name.count('*') + name.count('[')

    def createContext(self, cxxtype, cursor):
        ctx = TypeConvertContext()
        ctx.orig_type = cxxtype
        ctx.convable_type = self.TypeToConvertable(cxxtype)
        ctx.can_type = self.TypeToCanonical(cxxtype)
        ctx.const = self.TypeIsConst(cxxtype) or self.TypeIsConst(ctx.convable_type)

        ctx.orig_type_name = ctx.orig_type.spelling
        ctx.can_type_name = ctx.can_type.spelling
        ctx.convable_type_name = ctx.convable_type.spelling

        if ctx.const: ctx.can_type_name = self.TypeNameTrimConst(ctx.can_type_name)
        if self.TypeIsVolatile(ctx.can_type):
            ctx.can_type_name = self.TypeNameTrimConst(ctx.can_type_name)

        ctx.cursor = cursor
        return ctx

    def dumpContext(self, ctx, exit_ = True):
        tb = traceback.extract_stack()
        # print(tb, len(tb))
        lno = tb[len(tb) - 2][1]
        fn = tb[len(tb) - 2 ][2]
        posig = '%s:%s' % (fn, lno)

        print(posig, ctx.orig_type.spelling, ctx.orig_type.kind, ctx.orig_type_name, "\n",
              ctx.convable_type.kind, ctx.convable_type_name, "cva<-\n->can",
              ctx.can_type.kind, ctx.can_type_name,
              'const:', ctx.const, ctx.pointer_level, ctx.cursor.spelling, ctx.cursor.location)

        tdef = ctx.orig_type.get_declaration()
        if tdef is not None:
            print(tdef.kind, tdef.location, tdef.type.kind)
        if exit_: raise 'dumpctx traced.'
        return

    pass


class TypeConvForRust(TypeConv):
    tymap = {
        'bool': ['i8', 'c_char'], 'int': ['i32', 'c_int'], 'uint': ['u32', 'c_uint'],
        'unsigned int': ['u32', 'c_uint'],
        # 'long': ['i32', 'c_long'],  # wtf, 这个32位系统和64位系统不一样怎么办
        'long': ['i64', 'c_long'],  # wtf, 这个32位系统和64位系统不一样怎么办
        'unsigned long': ['u64', 'c_ulong'],
        'long long': ['i64', 'c_longlong'], 'unsigned long long': ['u64', 'c_ulonglong'],
        'short': ['i16', 'c_short'], 'unsigned short': ['u16', 'c_ushort'],
        'float': ['f32', 'c_float'], 'double': ['f64', 'c_double'],
        'char': ['i8', 'c_char'], 'unsigned char': ['u8', 'c_uchar'],
    }

    def __init__(self):
        super(TypeConvForRust, self).__init__()
        return

    # @param cxxtype clidx.Type
    # @return str
    # @return None 返回值为空，无返回值，不处理返回值
    def Type2RustArg(self, cxxtype):

        return

    # @param cxxtype clidx.Type
    # @return str
    # @return None 返回值为空，无返回值，不处理返回值
    def Type2RustRet(self, cxxtype, cursor):
        ctx = self.createContext(cxxtype, cursor)

        if ctx.convable_type.kind == clidx.TypeKind.POINTER:
            return self.Type2RustRetPointer(ctx)

        if ctx.convable_type.kind == clidx.TypeKind.LVALUEREFERENCE:
            return self.Type2RustRetLVRef(ctx)
        if ctx.convable_type.kind == clidx.TypeKind.RVALUEREFERENCE:
            return self.Type2RustRetLVRef(ctx)

        if ctx.convable_type.kind == clidx.TypeKind.VOID:
            return '()'

        if ctx.convable_type.kind == clidx.TypeKind.RECORD:
            return self.Type2RustRetRecord(ctx)

        if ctx.convable_type.kind == clidx.TypeKind.UNEXPOSED:
            return '(/*unexposed*/)'

        if ctx.convable_type.kind == clidx.TypeKind.ENUM:
            return 'i32'

        # 原始类型值类型
        if ctx.convable_type_name in TypeConvForRust.tymap:
            return self.Type2RustRetPrimitive(ctx)

        self.dumpContext(ctx)
        return ctx.orig_type.spelling

    # @param cxxtype clidx.Type
    # @return str
    def Type2ExtArg(self, cxxtype):
        return

    def Type2ExtRet(self, cxxtype):
        return

    def Type2RustRetFunctionProto(self, ctx):
        rety = '*mut c_void'
        rety = ctx.orig_type_name
        return rety

    def Type2RustRetPointer(self, ctx):
        ctx.pointer_level += 1
        pointee_type = ctx.convable_type.get_pointee()
        if pointee_type.kind == clidx.TypeKind.POINTER:
            ctx.pointer_level += 1
            pointee_type2 = pointee_type.get_pointee()
            if pointee_type2.kind == clidx.TypeKind.POINTER:
                ctx.pointer_level += 1
        if ctx.pointer_level > 2:
            glog.debug("")
            print('wtf, two many pointer level:' + str(ctx.pointer_level))
            exit(0)

        if ctx.can_type.kind == clidx.TypeKind.FUNCTIONPROTO:
            return self.Type2RustRetFunctionProto(ctx)

        if ctx.can_type.kind == clidx.TypeKind.RECORD:
            return ctx.can_type_name
        if ctx.can_type.kind == clidx.TypeKind.VOID:
            return '*mut c_void'

        can_name = ctx.can_type_name
        if ctx.const: can_name = self.TypeNameTrimConst(can_name)
        if can_name in TypeConvForRust.tymap:
            can_rsty = TypeConvForRust.tymap[can_name][0]
            if ctx.pointer_level == 1:
                if self.IsCharType(can_name): rety = 'String'
                else: rety = '*mut %s' % (can_rsty)
            elif ctx.pointer_level == 2:
                if self.IsCharType(can_name): rety = 'Vec<String>'
                else: rety = '*mut *mut %s' & (can_rsty)
            else: raise('not possible')
            return rety
        else:
            glog.debug("");
            print(678, 'wtf, type not in tymap:', can_name, ctx.orig_type_name, ctx.orig_type.spelling,
                  ctx.cursor.spelling, ctx.cursor.kind, ctx.cursor.semantic_parent.spelling)
            self.dumpContext(ctx)

        raise('not possible')
        return

    def Type2RustRetLVRef(self, ctx):
        # self.dumpContext(ctx)
        ctx.pointer_level += 1
        pointee_type = ctx.convable_type.get_pointee()
        # print(pointee_type.kind, pointee_type.spelling, ctx.can_type_name)

        if ctx.can_type.kind == clidx.TypeKind.RECORD:
            return ctx.can_type_name

        can_name = ctx.can_type_name
        if ctx.const: can_name = self.TypeNameTrimConst(can_name)
        if can_name in TypeConvForRust.tymap:
            can_rsty = TypeConvForRust.tymap[can_name][0]
            if self.IsCharType(can_name): rety = 'String'
            else: rety = '%s' % (can_rsty)
            return rety
        else:
            glog.debug("");
            print(678, 'wtf, type not in tymap:', can_name, ctx.orig_type_name, ctx.orig_type.spelling,
                  ctx.cursor.spelling, ctx.cursor.kind, ctx.cursor.semantic_parent.spelling)
            self.dumpContext(ctx)

        self.dumpContext(ctx)
        raise('not possible')

        return

    def Type2RustRetPrimitive(self, ctx):
        rety = TypeConvForRust.tymap[ctx.can_type_name][0]
        return rety

    def Type2RustRetRecord(self, ctx):
        rety = ctx.can_type_name
        return rety

    # TODO, like QList<int>
    def Type2RustRetUnexposed(self, ctx):
        rety = ctx.can_type_name
        return rety

    # @param cxxtype clidx.Type
    def TypeCXX2Rust(self, cxxtype, cursor, inty=False):
        # TODO c++ char => rust char
        raw_type_map = TypeConvForRust.tymap
        ctx = self.createContext(cxxtype, cursor)

        def is_const(ty): return ty.spelling.startswith('const ')

        can_type = self.TypeToCanonical(cxxtype)
        can_name = self.TypeCanName(can_type)
        can_name = ctx.can_type_name

        if cxxtype.kind == clidx.TypeKind.FUNCTIONPROTO:
            # TODO
            self.dumpContext(ctx)

        mut_or_no = 'mut'
        if is_const(cxxtype): mut_or_no = ''

        if cxxtype.kind == clidx.TypeKind.LVALUEREFERENCE:
            # TODO 为什么std::string &的can_type_name是i32呢？
            if ctx.can_type_name in raw_type_map:
                if self.IsCharType(ctx.can_type_name) or self.IsCharType(ctx.convable_type_name):
                    return "&%s String" % (mut_or_no)
                else:
                    return "&%s %s" % (mut_or_no, raw_type_map[ctx.can_type_name][0])

        if cxxtype.kind == clidx.TypeKind.POINTER or \
           cxxtype.kind == clidx.TypeKind.LVALUEREFERENCE:
            mut_or_no = 'mut'
            if is_const(cxxtype): mut_or_no = ''
            if 'wstring' in ctx.can_type_name:
                self.dumpContext(ctx)
            if ctx.can_type_name in raw_type_map:
                if self.IsCharType(ctx.can_type_name) or self.IsCharType(ctx.convable_type_name):
                    return "&%s String" % (mut_or_no)
                else:
                    return "&%s Vec<%s>" % (mut_or_no, raw_type_map[ctx.can_type_name][0])

            if ctx.can_type.kind == clidx.TypeKind.RECORD:
                if inty is True: return '&' + ctx.can_type_name
                else: return ctx.can_type_name
            if ctx.can_type.kind == clidx.TypeKind.FUNCTIONPROTO:
                return ctx.cursor.spelling
            if ctx.can_type.kind == clidx.TypeKind.VOID:
                return '*mut c_void'
            #if ctx.can_type.kind == clidx.TypeKind.UCHAR:
            #    return '& String'
            #if ctx.can_type.kind == clidx.TypeKind.CHAR:
            #    return '& String'
            if ctx.can_type.kind == clidx.TypeKind.CHAR16:
                return '& String'
            if ctx.can_type.kind == clidx.TypeKind.CHAR32:
                return '& String'
            if ctx.can_type.kind == clidx.TypeKind.WCHAR:
                return '& String'
            if ctx.can_type.kind == clidx.TypeKind.ENUM:
                return '&mut i32'
            print(can_name, ctx.can_type.kind)
            self.dumpContext(ctx)

        if cxxtype.kind == clidx.TypeKind.RVALUEREFERENCE:
            return "&mut %s" % (cxxtype.spelling.split(' ')[0])

        if cxxtype.kind == clidx.TypeKind.ENUM:
            return 'i32'

        if cxxtype.kind in [clidx.TypeKind.BOOL,
                            clidx.TypeKind.INT, clidx.TypeKind.UINT,
                            clidx.TypeKind.SHORT, clidx.TypeKind.USHORT,
                            clidx.TypeKind.LONG, clidx.TypeKind.ULONG,
                            clidx.TypeKind.LONGLONG, clidx.TypeKind.ULONGLONG,
                            clidx.TypeKind.CHAR_S,
                            clidx.TypeKind.UCHAR,
                            clidx.TypeKind.DOUBLE, clidx.TypeKind.FLOAT, ]:
            if ctx.can_type_name in raw_type_map:
                return '%s' % (raw_type_map[ctx.can_type_name][0])
            self.dumpContext(ctx)

        if cxxtype.spelling in ['uint']:
            raw_type_name = cxxtype.spelling
            if raw_type_name in raw_type_map:
                return '%s' % (raw_type_map[raw_type_name][0])
            self.dumpContext(ctx)

        if cxxtype.kind == clidx.TypeKind.TYPEDEF:
            under_type = cxxtype.get_declaration().underlying_typedef_type
            if under_type.kind == clidx.TypeKind.UNEXPOSED and \
               under_type.spelling.startswith('QFlags<'):
                return 'i32'
            return self.TypeCXX2Rust(under_type, cxxtype.get_declaration())

        # ### UNEXPOSED
        # maybe TODO，可能是python-clang绑定功能不全啊，模板解析不出来
        if cxxtype.kind == clidx.TypeKind.UNEXPOSED:
            if ctx.can_type.kind == clidx.TypeKind.ENUM:
                return 'i32'
            if cxxtype.get_declaration() is not None:
                tdef = cxxtype.get_declaration()
                if tdef.kind == clidx.CursorKind.CLASS_DECL:
                    return self.TypeCXX2Rust(tdef.type, tdef)
                if tdef.kind == clidx.CursorKind.STRUCT_DECL:
                    return self.TypeCXX2Rust(tdef.type, tdef)

            import re
            template_exp = '([a-zA-Z]+)\<([ a-zA-Z:]+)([\*])?\>'
            template_res = re.findall(template_exp, ctx.can_type_name)

            if ctx.can_type_name.startswith('Qt::'):
                return 'i32'
            if cxxtype.spelling == '::quintptr':
                return '*mut i32'
            if ctx.can_type_name.startswith('std::initializer_list'):
                return ctx.can_type_name.replace('std::initializer_list', 'QList')

            if ctx.can_type.kind == clidx.TypeKind.RECORD:
                if ctx.can_type_name.startswith('QFlags<'):
                    return 'i32'
                # (QList, QAction, *)
                # (QVector, unsigned int)
                # print(template_res, template_res[0], len(template_res[0]))
                if len(template_res) == 0:
                    self.dumpContext(ctx)
                if len(template_res[0]) == 2 or len(template_res[0]) == 3:
                    cls = template_res[0][0].strip()
                    inty = template_res[0][1].strip()
                    if cls in ['QVector', 'QList']: result_cls = 'Vec'
                    else: result_cls = cls
                    if inty in raw_type_map: result_inty = raw_type_map[inty][0]
                    else: result_inty = inty
                    return('%s<%s>' % (result_cls, result_inty))

            self.dumpContext(ctx)

        # TODO const char *const [], TypeKind.INCOMPLETEARRAY
        # 也有可能是一维INCOMPLETEARRAY, 所以还要考虑维度
        if cxxtype.kind == clidx.TypeKind.INCOMPLETEARRAY or \
           cxxtype.kind == clidx.TypeKind.CONSTANTARRAY:
            pointee_type = can_type.get_pointee()  # INVALID
            # print(666, can_type.kind, can_name, cxxtype.spelling)
            constq = 'mut'
            rsty = ''
            if self.TypeIsConst(cxxtype): constq = ''
            can_name = can_name.split(' ')[0]
            ext_name = can_name
            if can_name in raw_type_map: ext_name = raw_type_map[can_name][0]
            rsty = '&%s Vec<&%s %s>' % (constq, constq, ext_name)
            return rsty

        # TODO TypeKind.CONSTANTARRAY

        # ###### RECORD
        if cxxtype.kind == clidx.TypeKind.RECORD:
            return ctx.can_type_name

        if cxxtype.kind == clidx.TypeKind.MEMBERPOINTER:
            return '*mut u64'

        self.dumpContext(ctx)
        # glog.debug('just use default type name: ' + str(cxxtype.spelling) + ', ' + str(cxxtype.kind))
        return cxxtype.spelling

    def TypeCXX2RustExtern(self, cxxtype, cursor):
        raw_type_map = TypeConvForRust.tymap
        ctx = self.createContext(cxxtype, cursor)

        def is_const(ty): return ty.spelling.startswith('const ')

        can_type = self.TypeToCanonical(cxxtype)
        can_name = self.TypeCanName(can_type)
        can_name = ctx.can_type_name

        if cxxtype.kind == clidx.TypeKind.POINTER or \
           cxxtype.kind == clidx.TypeKind.LVALUEREFERENCE:
            mut_or_const = '*const ' if is_const(cxxtype) else '*mut '
            mut_or_const = '*mut '
            if can_name in raw_type_map:
                return mut_or_const + '%s' % (raw_type_map[can_name][1])
            if ctx.can_type.kind == clidx.TypeKind.RECORD:
                return mut_or_const + 'c_void'
            if ctx.can_type.kind == clidx.TypeKind.VOID:
                return mut_or_const + 'c_void'
            if ctx.can_type.kind == clidx.TypeKind.FUNCTIONPROTO:
                return mut_or_const + 'c_void'
            if ctx.can_type.kind == clidx.TypeKind.CHAR32:
                return mut_or_const + 'c_char'
            if ctx.can_type.kind == clidx.TypeKind.CHAR16:
                return mut_or_const + 'c_char'
            if ctx.can_type.kind == clidx.TypeKind.WCHAR:
                return mut_or_const + 'wchar_t'
            if ctx.can_type.kind == clidx.TypeKind.ENUM:
                return mut_or_const + 'c_int'
            self.dumpContext(ctx)

        if cxxtype.kind == clidx.TypeKind.RVALUEREFERENCE:
            return '*mut %s' % ('c_void')

        if cxxtype.kind == clidx.TypeKind.ENUM:
            return 'c_int'

        if cxxtype.kind in [clidx.TypeKind.BOOL,
                            clidx.TypeKind.INT, clidx.TypeKind.UINT,
                            clidx.TypeKind.SHORT, clidx.TypeKind.USHORT,
                            clidx.TypeKind.LONG, clidx.TypeKind.ULONG,
                            clidx.TypeKind.LONGLONG, clidx.TypeKind.ULONGLONG,
                            clidx.TypeKind.CHAR_S,
                            clidx.TypeKind.UCHAR,
                            clidx.TypeKind.DOUBLE, clidx.TypeKind.FLOAT, ]:
            if ctx.can_type_name in raw_type_map:
                return '%s' % (raw_type_map[ctx.can_type_name][1])
            self.dumpContext(ctx)

        if cxxtype.spelling in ['uint']:
            raw_type_name = cxxtype.spelling
            if raw_type_name in raw_type_map:
                return '%s' % (raw_type_map[raw_type_name][1])
            else:
                return '%s' % (raw_type_name)

        if cxxtype.kind == clidx.TypeKind.TYPEDEF:
            under_type = cxxtype.get_declaration().underlying_typedef_type
            if under_type.kind == clidx.TypeKind.UNEXPOSED and \
               under_type.spelling.startswith('QFlags<'):
                return 'c_int'
            return self.TypeCXX2RustExtern(under_type, cxxtype.get_declaration())

        # maybe TODO
        if cxxtype.kind == clidx.TypeKind.UNEXPOSED:
            if ctx.can_type.kind == clidx.TypeKind.ENUM:
                return 'c_int'
            if ctx.can_type.kind == clidx.TypeKind.UINT:
                return 'c_int'
            if cxxtype.get_declaration() is not None:
                tdef = cxxtype.get_declaration()
                if tdef.kind == clidx.CursorKind.TYPEDEF_DECL:
                    return self.TypeCXX2RustExtern(tdef.type, tdef)
                if tdef.kind == clidx.CursorKind.CLASS_DECL:
                    return self.TypeCXX2RustExtern(tdef.type, tdef)
                if tdef.kind == clidx.CursorKind.STRUCT_DECL:
                    return self.TypeCXX2RustExtern(tdef.type, tdef)

            import re
            template_exp = '([a-zA-Z]+)\<([ a-zA-Z:]+)([\*])?\>'
            template_res = re.findall(template_exp, ctx.can_type_name)

            # like QAccessible::State
            if ctx.can_type_name.startswith('Qt::') or (
                    ctx.can_type_name[0] == 'Q' and '::' in ctx.can_type_name):
                return 'c_int'
            if cxxtype.spelling == '::quintptr':
                return '*mut c_uint'
            # under_type = cxxtype.get_declaration().underlying_typedef_type
            # print(777, under_type.spelling, under_type.kind)

            if cxxtype.spelling.startswith('std::initializer_list'):
                return ctx.can_type_name.replace('std::initializer_list', 'QList')

            # TODO fix
            # QPair<T1, T2>, QMap<T1, T2>, QHash<T1, T2>
            # QGenericMatrix<3, 3, float>
            if ctx.can_type.kind == clidx.TypeKind.RECORD:
                if 'QPair<' in ctx.can_type_name:
                    return ctx.can_type_name
                if 'QMap<' in ctx.can_type_name:
                    return ctx.can_type_name
                if 'QHash<' in ctx.can_type_name:
                    return ctx.can_type_name
                if 'QGenericMatrix<' in ctx.can_type_name:
                    return ctx.can_type_name
                if 'QFlags<' in ctx.can_type_name:
                    return ctx.can_type_name

            if ctx.can_type.kind == clidx.TypeKind.RECORD:
                # (QList, QAction, *)
                # (QVector, unsigned int)
                print(template_res, ctx.can_type_name)
                if len(template_res) == 0:
                    self.dumpContext(ctx)
                # print(template_res, template_res[0], len(template_res[0]))
                if len(template_res[0]) == 2 or len(template_res[0]) == 3:
                    cls = template_res[0][0].strip()
                    inty = template_res[0][1].strip()
                    if cls in ['QVector', 'QList']: result_cls = 'Vec'
                    else: result_cls = cls
                    if inty in raw_type_map: result_inty = raw_type_map[inty][0]
                    else: result_inty = inty
                    return('%s<%s>' % (result_cls, result_inty))

            if 'std::string' in ctx.convable_type_name:
                return ctx.convable_type_name
            if 'std::wstring' in ctx.convable_type_name:
                return ctx.convable_type_name
            if 'std::u16string' in ctx.convable_type_name:
                return ctx.convable_type_name
            if 'std::u32string' in ctx.convable_type_name:
                return ctx.convable_type_name

            # for template T
            if 'type-parameter-' in ctx.can_type_name:
                self.dumpContext(ctx)

            self.dumpContext(ctx)

        if cxxtype.kind == clidx.TypeKind.RECORD:
            if is_const(cxxtype):
                # return '*const %s' % ('c_void')  # 不好处理，全换mut试试吧
                return '*mut %s' % ('c_void')
            else:
                return '*mut %s' % ('c_void')

        # TODO const char *const [], TypeKind.INCOMPLETEARRAY
        # 也有可能是一维INCOMPLETEARRAY, 所以还要考虑维度
        if cxxtype.kind == clidx.TypeKind.INCOMPLETEARRAY or \
            cxxtype.kind == clidx.TypeKind.CONSTANTARRAY:
            pointee_type = can_type.get_pointee()  # INVALID
            # print(666, can_type.kind, can_name, cxxtype.spelling)
            constq = 'mut'
            rsty = ''
            # if self.TypeIsConst(cxxtype): constq = 'const'
            can_name = can_name.split(' ')[0]
            ext_name = can_name
            if can_name in raw_type_map: ext_name = raw_type_map[can_name][1]
            rsty = '*%s *%s %s' % (constq, constq, ext_name)

            return rsty
            # print(666, rsty, cxxtype.spelling)
            # exit(0)
            pass

        if cxxtype.kind == clidx.TypeKind.VOID:
            return 'void'

        if cxxtype.kind == clidx.TypeKind.MEMBERPOINTER:
            return '*mut c_void'

        self.dumpContext(ctx)
        # print(888, 'just use default type name:', cxxtype.spelling, cxxtype.kind)
        return cxxtype.spelling

    @staticmethod
    def TypeCXXCanonical(cxxtype):

        clean_type = cxxtype.get_pointee()
        # print(444, canical_type.spelling, canical_type.spelling, canical_type.kind)
        def is_const(ty): return ty.spelling.startswith('const ')
        def trim_const(ty):
            tysegs = ty.spelling.split(' ') + ['const']
            tysegs.remove('const')
            return tysegs[0]
        if clean_type.kind == clidx.TypeKind.RECORD:
            tyname = trim_const(clean_type)
            return tyname
        return clean_type.spelling



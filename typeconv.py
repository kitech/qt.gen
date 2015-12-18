# encoding: utf8

import logging

import clang
import clang.cindex
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

        self.orig_cursor = None
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
            under_type = cxxtype.get_declaration().under_type()
            return self.TypeToCanonical(under_type)

        return cxxtype

    def TypeToConvertable(self, cxxtype):
        if cxxtype.kind == clidx.TypeKind.TYPEDEF:
            under_type = cxxtype.get_declaration().underlying_typedef_type
            return self.TypeToConvertable(under_type)
        return cxxtype

    def TypeCanName(self, cxxtype):
        canty = self.TypeToCanonical(cxxtype)
        if canty.kind == clidx.TypeKind.FUNCTIONPROTO:
            # print(444, "function proto:", canty.spelling)
            # exit(0)
            pass
        canname = canty.spelling
        if self.TypeIsConst(canty): canname = self.TypeTrimConst(canty)
        canname = canname.replace('unsigned ', 'u')
        return canname

    def TypeIsConst(self, cxxtype):
        return cxxtype.spelling.startswith('const ')

    def TypeTrimConst(self, cxxtype):
        tysegs = cxxtype.spelling.split(' ')
        if 'const' in tysegs: tysegs.remove('const')
        return ' '.join(tysegs)

    def TypeNameTrimConst(self, tyname):
        tysegs = tyname.split(' ')
        if 'const' in tysegs: tysegs.remove('const')
        return ' '.join(tysegs)

    def IsCharType(self, tyname): return 'char' in tyname

    def IsPointer(self, cxxtype):
        if cxxtype.kind == clidx.TypeKind.POINTER or \
           cxxtype.kind == clidx.TypeKind.LVALUEREFERENCE or \
           cxxtype.kind == clidx.TypeKind.RVALUEREFERENCE:
            return True
        return False

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

        ctx.orig_cursor = cursor
        return ctx

    def dumpContext(self, ctx):
        print(890, ctx.orig_type.kind, ctx.orig_type_name, "\n",
              ctx.convable_type.kind, ctx.convable_type_name, "cva<-\n->can",
              ctx.can_type.kind, ctx.can_type_name,
              ctx.const, ctx.pointer_level)
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


    # @param cxxtype clang.cindex.Type
    # @return str
    # @return None 返回值为空，无返回值，不处理返回值
    def Type2RustArg(self, cxxtype):

        return

    # @param cxxtype clang.cindex.Type
    # @return str
    # @return None 返回值为空，无返回值，不处理返回值
    def Type2RustRet(self, cxxtype, cursor):
        ctx = self.createContext(cxxtype, cursor)
        if 'QMetaObject' in ctx.orig_type_name:
            self.dumpContext(ctx)
            exit(0)

        if ctx.convable_type.kind == clidx.TypeKind.POINTER:
            return self.Type2RustRetPointer(ctx)

        if ctx.convable_type.kind == clidx.TypeKind.LVALUEREFERENCE:
            return self.Type2RustRetLVRef(ctx)

        if ctx.convable_type.kind == clidx.TypeKind.VOID:
            return '()'

        if ctx.convable_type.kind == clidx.TypeKind.RECORD:
            return self.Type2RustRetRecord(ctx)

        # 原始类型值类型
        if ctx.convable_type_name in TypeConvForRust.tymap:
            return self.Type2RustRetPrimitive(ctx)

        print(783, 'wtf')
        self.dumpContext(ctx)
        exit(0)
        return ctx.orig_type.spelling

    # @param cxxtype clang.cindex.Type
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
            if ctx.can_type.kind == clidx.TypeKind.RECORD:
                rety = can_name
                return rety
            if ctx.can_type.kind == clidx.TypeKind.VOID:
                rety = '*mut c_void'
                return rety
            glog.debug("");
            print(678, 'wtf, type not in tymap:', can_name, ctx.orig_type_name, ctx.orig_type.spelling,
                  ctx.orig_cursor.spelling, ctx.orig_cursor.kind, ctx.orig_cursor.semantic_parent.spelling)
            self.dumpContext(ctx)
            exit(0)

        raise('not possible')
        return

    def Type2RustRetLVRef(self, ctx):
        self.dumpContext(ctx)
        ctx.pointer_level += 1
        pointee_type = ctx.convable_type.get_pointee()
        print(pointee_type.kind, pointee_type.spelling, ctx.can_type_name)

        can_name = ctx.can_type_name
        if ctx.const: can_name = self.TypeNameTrimConst(can_name)
        if can_name in TypeConvForRust.tymap:
            can_rsty = TypeConvForRust.tymap[can_name][0]
            if self.IsCharType(can_name): rety = 'String'
            else: rety = '%s' % (can_rsty)
            return rety
        else:
            if ctx.can_type.kind == clidx.TypeKind.RECORD:
                rety = can_name
                return rety
            glog.debug("");
            print(678, 'wtf, type not in tymap:', can_name, ctx.orig_type_name, ctx.orig_type.spelling,
                  ctx.orig_cursor.spelling, ctx.orig_cursor.kind, ctx.orig_cursor.semantic_parent.spelling)
            self.dumpContext(ctx)
            exit(0)

        raise('not possible')

        exit(0)
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

    # @param cxxtype clang.cindex.Type
    def TypeCXX2Rust(self, cxxtype):
        # TODO c++ char => rust char
        raw_type_map = {
            'bool': 'i8', 'int': 'i32', 'uint': 'u32', 'unsigned int': 'u32',
            'long': 'i32', 'unsigned long': 'u32', 'long long': 'i64', 'unsigned long long': 'u64',
            'short': 'i16', 'ushort': 'u16', 'unsigned short': 'u16',
            'char': 'i8', 'uchar': 'u8', 'unsigned char': 'u8',
            'char16_t': 'i16', 'wchar_t': 'wchar_t',
            'char32_t': 'i32',
            'float': 'f32', 'double': 'f64',
            'qint64': 'i64',
            'qreal': 'f64',
            # very special
            'void': 'u8', 'std::string': 'u8', 'std::wstring': 'u16',
            'std::u32string': 'u32', 'std::u16string': 'u16',
        }
        # print(888, cxxtype.kind, cxxtype.spelling)

        def is_const(ty): return ty.spelling.startswith('const ')

        can_type = self.TypeToCanonical(cxxtype)
        can_name = self.TypeCanName(can_type)

        if cxxtype.kind == clang.cindex.TypeKind.FUNCTIONPROTO:
            # TODO
            glog.debug('function proto type')
            pass

        if cxxtype.kind == clang.cindex.TypeKind.POINTER or \
           cxxtype.kind == clang.cindex.TypeKind.LVALUEREFERENCE:
            mut_or_no = 'mut'
            if is_const(cxxtype): mut_or_no = ''

            if can_name in raw_type_map:
                if self.IsCharType(can_name):
                    return "&%s String" % (mut_or_no)
                else:
                    return "&%s %s" % (mut_or_no, raw_type_map[can_name])
            else:
                return "&%s %s" % (mut_or_no, can_name)

        if cxxtype.kind == clang.cindex.TypeKind.RVALUEREFERENCE:
            return "&mut %s" % (cxxtype.spelling.split(' ')[0])

        if cxxtype.kind == clang.cindex.TypeKind.ENUM:
            return 'i32'

        if cxxtype.kind in [clang.cindex.TypeKind.BOOL,
                            clang.cindex.TypeKind.INT, clang.cindex.TypeKind.UINT,
                            clang.cindex.TypeKind.SHORT, clang.cindex.TypeKind.USHORT,
                            clang.cindex.TypeKind.LONG, clang.cindex.TypeKind.ULONG,
                            clang.cindex.TypeKind.LONGLONG, clang.cindex.TypeKind.ULONGLONG,
                            clang.cindex.TypeKind.CHAR_S,
                            clang.cindex.TypeKind.UCHAR,
                            clang.cindex.TypeKind.DOUBLE, clang.cindex.TypeKind.FLOAT, ]:
            raw_type_name = cxxtype.spelling
            if raw_type_name in raw_type_map:
                return '%s' % (raw_type_map[raw_type_name])
            else:
                glog.debug('just use default type name: %s, %s', cxxtype.spelling, str(cxxtype.kind))
                if cxxtype.spelling == 'int':
                    exit(0)
                return '%s' % (self.TypeNameTrimConst(raw_type_name))

        if cxxtype.spelling in ['uint']:
            raw_type_name = cxxtype.spelling
            if raw_type_name in raw_type_map:
                return '%s' % (raw_type_map[raw_type_name])
            else:
                return '%s' % (raw_type_name)

        if cxxtype.kind == clang.cindex.TypeKind.TYPEDEF:
            under_type = cxxtype.get_declaration().underlying_typedef_type
            if under_type.kind == clang.cindex.TypeKind.UNEXPOSED and \
               under_type.spelling.startswith('QFlags<'):
                return 'i32'
            return self.TypeCXX2Rust(under_type)

        # maybe TODO
        if cxxtype.kind == clang.cindex.TypeKind.UNEXPOSED:
            if cxxtype.spelling.startswith('Qt::') or \
               (cxxtype.spelling.startswith('Q') and '::' in cxxtype.spelling):
                return 'i32'
            if cxxtype.spelling == '::quintptr':
                return '*mut i32'

        # TODO const char *const [], TypeKind.INCOMPLETEARRAY
        # 也有可能是一维INCOMPLETEARRAY, 所以还要考虑维度
        if cxxtype.kind == clang.cindex.TypeKind.INCOMPLETEARRAY or \
           cxxtype.kind == clang.cindex.TypeKind.CONSTANTARRAY:
            pointee_type = can_type.get_pointee()  # INVALID
            # print(666, can_type.kind, can_name, cxxtype.spelling)
            constq = 'mut'
            rsty = ''
            if self.TypeIsConst(cxxtype): constq = ''
            can_name = can_name.split(' ')[0]
            ext_name = can_name
            if can_name in raw_type_map: ext_name = raw_type_map[can_name]
            rsty = '&%s Vec<&%s %s>' % (constq, constq, ext_name)
            return rsty
            # print(666, rsty, cxxtype.spelling)
            # exit(0)
            pass
        # TODO TypeKind.CONSTANTARRAY

        glog.debug('just use default type name: ' + str(cxxtype.spelling) + ', ' + str(cxxtype.kind))
        return cxxtype.spelling

    def TypeCXX2RustExtern(self, cxxtype):
        raw_type_map = {
            'bool': 'int8_t', 'int': 'c_int', 'uint': 'c_uint', 'unsigned int': 'c_uint',
            'long': 'c_long', 'unsigned long': 'c_ulong', 'long long': 'c_longlong',
            'unsigned long long': 'uint64_t',
            'short': 'c_short', 'ushort': 'c_ushort', 'unsigned short': 'c_ushort',
            'char': 'c_char', 'uchar': 'c_uchar', 'unsigned char': 'c_uchar',
            'char16_t': 'int16_t', 'wchar_t': 'wchar_t',
            'char32_t': 'int32_t',
            'float': 'c_float', 'double': 'c_double',
            'qint64': 'int64_t',
            'qreal': 'c_double',
            # very special
            # 'void': 'c_void',
            'void': 'uint8_t',
            'std::string': 'c_char', 'std::wstring': 'c_char',
            'std::u32string': 'c_char', 'std::u16string': 'c_char',
        }

        def is_const(ty): return ty.spelling.startswith('const ')

        can_type = self.TypeToCanonical(cxxtype)
        can_name = self.TypeCanName(can_type)

        if cxxtype.kind == clang.cindex.TypeKind.POINTER or \
           cxxtype.kind == clang.cindex.TypeKind.LVALUEREFERENCE:
            if is_const(cxxtype):
                if can_name in raw_type_map:
                    return '*const %s' % (raw_type_map[can_name])
                else:
                    # return '*const %s' % ('c_void')  # 不好处理，全换mut试试吧
                    return '*mut %s' % ('c_void')
            else:
                if can_type.kind == clang.cindex.TypeKind.FUNCTIONPROTO:
                    pass  # TODO
                if can_name in raw_type_map:
                    return '*mut %s' % (raw_type_map[can_name])
                else:
                    return '*mut %s' % ('c_void')

        if cxxtype.kind == clang.cindex.TypeKind.RVALUEREFERENCE:
            return '*mut %s' % (cxxtype.spelling.split(' ')[0])

        if cxxtype.kind == clang.cindex.TypeKind.ENUM:
            return 'c_int'

        if cxxtype.kind in [clang.cindex.TypeKind.BOOL,
                            clang.cindex.TypeKind.INT, clang.cindex.TypeKind.UINT,
                            clang.cindex.TypeKind.SHORT, clang.cindex.TypeKind.USHORT,
                            clang.cindex.TypeKind.LONG, clang.cindex.TypeKind.ULONG,
                            clang.cindex.TypeKind.LONGLONG, clang.cindex.TypeKind.ULONGLONG,
                            clang.cindex.TypeKind.CHAR_S,
                            clang.cindex.TypeKind.UCHAR,
                            clang.cindex.TypeKind.DOUBLE, clang.cindex.TypeKind.FLOAT, ]:
            raw_type_name = cxxtype.spelling
            if raw_type_name in raw_type_map:
                return '%s' % (raw_type_map[raw_type_name])
            else:
                # print(888, 'just use default type name:', cxxtype.spelling, cxxtype.kind)
                if cxxtype.spelling == 'int':
                    exit(0)
                return '%s' % (self.TypeNameTrimConst(raw_type_name))

        if cxxtype.spelling in ['uint']:
            raw_type_name = cxxtype.spelling
            if raw_type_name in raw_type_map:
                return '%s' % (raw_type_map[raw_type_name])
            else:
                return '%s' % (raw_type_name)

        if cxxtype.kind == clang.cindex.TypeKind.TYPEDEF:
            under_type = cxxtype.get_declaration().underlying_typedef_type
            if under_type.kind == clang.cindex.TypeKind.UNEXPOSED and \
               under_type.spelling.startswith('QFlags<'):
                return 'c_int'
            return self.TypeCXX2RustExtern(under_type)

        # maybe TODO
        if cxxtype.kind == clang.cindex.TypeKind.UNEXPOSED:
            if cxxtype.spelling.startswith('Qt::') or \
               (cxxtype.spelling.startswith('Q') and '::' in cxxtype.spelling):
                return 'c_int'
            if cxxtype.spelling == '::quintptr':
                return '*mut c_uint'
            # under_type = cxxtype.get_declaration().underlying_typedef_type
            # print(777, under_type.spelling, under_type.kind)

        if cxxtype.kind == clang.cindex.TypeKind.RECORD:
            if is_const(cxxtype):
                # return '*const %s' % ('c_void')  # 不好处理，全换mut试试吧
                return '*mut %s' % ('c_void')
            else:
                return '*mut %s' % ('c_void')

        # TODO const char *const [], TypeKind.INCOMPLETEARRAY
        # 也有可能是一维INCOMPLETEARRAY, 所以还要考虑维度
        if cxxtype.kind == clang.cindex.TypeKind.INCOMPLETEARRAY or \
            cxxtype.kind == clang.cindex.TypeKind.CONSTANTARRAY:
            pointee_type = can_type.get_pointee()  # INVALID
            # print(666, can_type.kind, can_name, cxxtype.spelling)
            constq = 'mut'
            rsty = ''
            # if self.TypeIsConst(cxxtype): constq = 'const'
            can_name = can_name.split(' ')[0]
            ext_name = can_name
            if can_name in raw_type_map: ext_name = raw_type_map[can_name]
            rsty = '*%s *%s %s' % (constq, constq, ext_name)

            return rsty
            # print(666, rsty, cxxtype.spelling)
            # exit(0)
            pass


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
        if clean_type.kind == clang.cindex.TypeKind.RECORD:
            tyname = trim_const(clean_type)
            return tyname
        return clean_type.spelling



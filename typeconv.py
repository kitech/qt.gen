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


# for python2 base class need to object
class TypeConv(object):
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
    pass


class TypeConvForRust(TypeConv):

    def __init__(self):
        super(TypeConvForRust, self).__init__()
        return

    # @param cxxtype clang.cindex.Type
    def TypeCXX2Rust(self, cxxtype):
        raw_type_map = {
            'bool': 'i8', 'int': 'i32', 'uint': 'u32',
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

        glog.debug('just use default type name: ' + str(cxxtype.spelling) + ', ' + str(cxxtype.kind))
        return cxxtype.spelling

    def TypeCXX2RustExtern(self, cxxtype):
        raw_type_map = {
            'bool': 'int8_t', 'int': 'c_int', 'uint': 'c_uint',
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
        # print(888, cxxtype.kind, cxxtype.spelling)

        def is_const(ty): return ty.spelling.startswith('const ')

        can_type = self.TypeToCanonical(cxxtype)
        can_name = self.TypeCanName(can_type)

        if cxxtype.kind == clang.cindex.TypeKind.POINTER or \
           cxxtype.kind == clang.cindex.TypeKind.LVALUEREFERENCE:
            if is_const(cxxtype):
                if can_name in raw_type_map:
                    return '*const %s' % (raw_type_map[can_name])
                else:
                    return '*const %s' % ('c_void')
            else:
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
                print(888, 'just use default type name:', cxxtype.spelling, cxxtype.kind)
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

        if cxxtype.kind == clang.cindex.TypeKind.RECORD:
            if is_const(cxxtype):
                return '*const %s' % ('c_void')
            else:
                return '*mut %s' % ('c_void')

        print(888, 'just use default type name:', cxxtype.spelling, cxxtype.kind)
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



package main

import (
	"encoding/hex"
	"fmt"
	"gopp"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-clang/v3.9/clang"
	funk "github.com/thoas/go-funk"
)

func get_decl_loc(cursor clang.Cursor) string {
	loc := cursor.Location()
	file, lineno, _, _ := loc.FileLocation()
	// log.Println(file.Name())
	return fmt.Sprintf("%s:%d", file.Name(), lineno)
}

// spell location
func get_decl_sploc(cursor clang.Cursor) string {
	loc := cursor.Location()
	file, lineno, _, _ := loc.SpellingLocation()
	// log.Println(file.Name())
	return fmt.Sprintf("%s:%d", file.Name(), lineno)
}

// xxx/include/ => /usr/include/qt
func fix_inc_name(name string) string {
	if !strings.HasPrefix(name, "/usr/include/qt") {
		if strings.Contains(name, "/include/") {
			return "/usr/include/qt/" + strings.Split(name, "/include/")[1]
		}
	}
	return name
}

// # like core without qt prefix
func get_decl_mod_lower(cursor clang.Cursor) string { return get_decl_mod(cursor) }
func get_decl_mod(cursor clang.Cursor) string {
	loc := cursor.Location()
	file, _, _, _ := loc.FileLocation()
	log.Println(cursor.Spelling(), cursor.IsCursorDefinition(), file.Name())
	if !strings.HasPrefix(file.Name(), "/usr/include/qt") &&
		!strings.Contains(file.Name(), "/gcc_64/include/") {
		if strings.Contains(file.Name(), "bsheaders/QtCore/") { // fix qRegisterResourceData
			return "core"
		}
		return "stdglobal"
	}
	segs := strings.Split(file.Name(), "/")
	// log.Println(segs[len(segs)-2], segs[len(segs)-2][2:])
	dmod := strings.ToLower(segs[len(segs)-2][2:])
	if !strings.HasPrefix(filepath.Base(filepath.Dir(file.Name())), "Qt") {
		dmod = filepath.Base(filepath.Dir(file.Name()))
	}
	log.Println(cursor.Spelling(), dmod, file.Name())
	if dmod == "android" || dmod == "jni" {
		dmod = "androidextras"
	}
	log.Println(cursor.Spelling(), dmod, file.Name())
	if _, found := modDeps[dmod]; !found {
		if false {
			log.Printf("unknown module: %s, %s %s, %s\n",
				dmod, cursor.Spelling(), file.Name(), filepath.Base(file.Name()))
			time.Sleep(500 * time.Millisecond)
		}
	}

	return dmod
}

// camal case
func get_decl_mod_norm(cursor clang.Cursor) string { return get_decl_mod(cursor) }

// 计算包名补全
func calc_package_prefix(curc, refc clang.Cursor) string {
	curmod := get_decl_mod(curc)
	refmod := get_decl_mod(refc)
	if refmod != curmod {
		return "qt" + refmod + "."
	}
	return ""
}

func is_qt_class(ty clang.Type) bool {
	nty := get_bare_type(ty)
	name := nty.Spelling()
	if len(name) < 2 {
		return false
	}
	// QImageCleanupFunction
	if name[0:1] == "Q" && strings.ToUpper(name[1:2]) == name[1:2] &&
		!strings.Contains(name, "::") && !strings.HasSuffix(name, "Function") {
		return true
	}
	return false
}

func is_private_method(c clang.Cursor) bool {
	return c.Kind() == clang.Cursor_CXXMethod &&
		c.AccessSpecifier() == clang.AccessSpecifier_Private
}

// 去掉reference和pointer,并查找其定义类型名，不带const
// QCameraFocusZoneList => QList???
func get_bare_type(ty clang.Type) clang.Type {
	switch ty.Kind() {
	case clang.Type_LValueReference, clang.Type_Pointer:
		if ty.PointeeType().Kind() == clang.Type_Void {
			// return ty
		}
		return get_bare_type(ty.PointeeType())
	}

	return ty.Declaration().Type()
}

func is_go_keyword(s string) bool {
	keywords := map[string]int{"match": 1, "type": 1, "move": 1, "select": 1, "case": 1,
		"map": 1, "range": 1, "var": 1, "len": 1, "fmt": 1, "err": 1, "go": 1, "func": 1,
		"package": 1, "import": 1,
		"begin": 1, "end": 1,
		"out": 1, "include": 1, "extern": 1, "module": 1, "require": 1}
	_, ok := keywords[s]
	return ok
}

func is_rs_keyword(s string) bool {
	keywords := map[string]int{"match": 1, "type": 1, "move": 1, "select": 1, "case": 1,
		"map": 1, "range": 1, "var": 1, "len": 1, "fmt": 1, "err": 1,
		"pub": 1, "fn": 1, "mod": 1, "use": 1, "self": 1, "super": 1}
	_, ok := keywords[s]
	return ok
}

// 包含1个以上的纯虚方法
// TODO 父类也有纯虚方法，并且当前类没有实现该方法
func is_pure_virtual_class(cursor clang.Cursor) bool {
	// pure virtual class check
	pure_virtual_class := false
	extraPureClses := map[string]int{"QAnimationGroup": 1, "QAccessibleObject": 1}
	if _, ok := extraPureClses[cursor.Spelling()]; ok {
		return true
	}
	if strings.HasPrefix(cursor.Spelling(), "QAbstract") {
		return true
	}
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.CXXMethod_IsPureVirtual() {
			pure_virtual_class = true
			return clang.ChildVisit_Break
		}
		return clang.ChildVisit_Continue
	})
	return pure_virtual_class
}

func find_base_classes(cursor clang.Cursor) []clang.Cursor {
	bcs := make([]clang.Cursor, 0)
	cursor.Visit(func(c, p clang.Cursor) clang.ChildVisitResult {
		if c.Kind() == clang.Cursor_CXXBaseSpecifier {
			if _, ok := privClasses[c.Definition().Spelling()]; !ok {
				bcs = append(bcs, c.Definition())
			}
		}
		if c.Kind() == clang.Cursor_CXXMethod {
			return clang.ChildVisit_Break
		}
		return clang.ChildVisit_Continue
	})
	return bcs
}

func has_qobject_base_class(cursor clang.Cursor) bool {
	bcs := find_base_classes(cursor)
	for _, bc := range bcs {
		if bc.Spelling() == "QObject" {
			return true
		}
		if has_qobject_base_class(bc) {
			return true
		}
	}
	return false
}

// for class
func has_virtual_projected(cursor clang.Cursor, include_base bool) (has bool) {
	cursor.Visit(func(c, p clang.Cursor) clang.ChildVisitResult {
		switch c.Kind() {
		case clang.Cursor_Constructor:
		case clang.Cursor_CXXMethod:
			if c.AccessSpecifier() == clang.AccessSpecifier_Protected &&
				(c.CXXMethod_IsPureVirtual() || c.CXXMethod_IsPureVirtual()) {
				has = true
			}
		}
		if has {
			return clang.ChildVisit_Break
		}
		return clang.ChildVisit_Continue
	})
	if has {
		return
	}
	if include_base {
		bcs := find_base_classes(cursor)
		if len(bcs) > 0 {
			return has_virtual_projected(bcs[0], include_base)
		}
	}
	return
}

func has_copy_ctor(cursor clang.Cursor) bool {
	has := false
	cursor.Visit(func(c, p clang.Cursor) clang.ChildVisitResult {
		switch c.Kind() {
		case clang.Cursor_Constructor:
			if c.AccessSpecifier() == clang.AccessSpecifier_Public &&
				(c.CXXConstructor_IsCopyConstructor() ||
					c.CXXConstructor_IsMoveConstructor()) {
				has = true
				return clang.ChildVisit_Break
			}
		case clang.Cursor_CXXMethod:

		}
		return clang.ChildVisit_Continue
	})
	return has
}

// has default ctor, no virtual method, no virtual base class
func is_trivial_class(cursor clang.Cursor) bool {
	hasDefaultCtor := false
	hasVirtMethod := false
	hasVirtBaseCls := false
	cursor.Visit(func(c, p clang.Cursor) clang.ChildVisitResult {
		switch c.Kind() {
		case clang.Cursor_Constructor:
			if c.AccessSpecifier() == clang.AccessSpecifier_Public &&
				c.CXXConstructor_IsDefaultConstructor() {
				hasDefaultCtor = true
			}
		case clang.Cursor_CXXMethod:
			hasVirtMethod = hasVirtMethod || c.CXXMethod_IsVirtual()
		case clang.Cursor_CXXBaseSpecifier:
			hasVirtBaseCls = hasVirtBaseCls || c.IsVirtualBase()
		}
		return clang.ChildVisit_Recurse
	})
	return hasDefaultCtor && !hasVirtMethod && !hasVirtBaseCls
}

// TODO auto detect
// non explicit ctor by user. ctor is not public
func is_deleted_class(cursor clang.Cursor) bool {
	deleted := false
	arr := map[string]int{"QClipboard": 1, "QInputMethod": 1, "QSessionManager": 1,
		"QPaintDevice": 1, "QPagedPaintDevice": 1, "QScroller": 1, "QStandardPaths": 1,
		"QLoggingCategory": 1, "QCameraExposure": 1, "QCameraFocus": 1, "QCameraImageProcessing": 1,
		"QOpenGLPaintDevice": 1}
	if _, ok := arr[cursor.Spelling()]; ok {
		return true
	}
	// go-clang not detect deleted feature now.
	cursor.Visit(func(c, p clang.Cursor) clang.ChildVisitResult {
		switch c.Kind() {
		case clang.Cursor_Destructor:
		}
		return clang.ChildVisit_Recurse
	})
	return deleted
}

var _cmgl = NewIncMangler()

// TODO auto detect
// like end with: = delete;
func is_deleted_method(cursor, parent clang.Cursor) bool {
	mths := map[string]int{
		"_ZN18QAbstractItemModelaSERKS_": 1, "_ZN18QAbstractAnimationaSERKS_": 1,
		"_ZN15QAnimationGroupaSERKS_": 1, "_ZN18QAbstractListModelaSERKS_": 1,
		"_ZN26QAbstractNativeEventFilteraSERKS_": 1, "_ZN19QAbstractProxyModelaSERKS_": 1,
		"_ZN14QAbstractStateaSERKS_": 1, "_ZN19QAbstractTableModelaSERKS_": 1,
		"_ZN19QAbstractTransitionaSERKS_": 1, "_ZN7QBufferaSERKS_": 1,
		"_ZN18QCommandLineParseraSERKS_": 1, "_ZNK16QLoggingCategoryclEv": 1,
		"_ZN11QDataStreamrsER8qfloat16": 1, // because of qfloat16 class
		"_ZN11QDataStreamlsE8qfloat16":  1,
		"_ZN18QSemaphoreReleaserC2EOS_": 1, "_ZN18QSemaphoreReleaseraSEOS_": 1,
		"_ZN14QSignalBlockeraSEOS_": 1, "_ZN14QSignalBlockerC2EOS_": 1,
		"_ZN18QOpenGLPaintDeviceaSERKS_": 1, "_ZN15QGraphicsObjectC2EP13QGraphicsItem": 1,

		"_ZN16QOpenGLFunctions14glShaderSourceEjiPPKcPKi":                              1,
		"_ZNK23QOperatingSystemVersion11isAnyOfTypeESt16initializer_listINS_6OSTypeEE": 1,
	}
	mname := _cmgl.origin(cursor)
	if _, ok := mths[mname]; ok {
		return true
	}
	return false
}

func is_pure_virtual_method(cursor, parent clang.Cursor) bool {
	return cursor.CXXMethod_IsPureVirtual()
}

// copy ctor is deleted
func is_nocopy_class(cursor clang.Cursor) bool {
	arr := map[string]int{"AbstractComparatorFunction": 1}
	if _, ok := arr[cursor.Spelling()]; ok {
		return true
	}
	return false
}

func getOverloadedIndex(cursor clang.Cursor, cursors []clang.Cursor) int {
	idx := 0
	for _, cs := range cursors {
		if cs.Kind() == cursor.Kind() && cs.Spelling() == cursor.Spelling() {
			if cs.Equal(cursor) {
				return idx
			} else {
				idx += 1
			}
		}
	}
	return idx
}

// TODO auto detect
func is_qt_private_class(cursor clang.Cursor) bool {
	loc := cursor.Type().Declaration().Definition().Location()
	file, _, _, _ := loc.FileLocation()
	// log.Println(file.Name(), cursor.Spelling(), cursor.IsCursorDefinition(), cursor.Definition().IsCursorDefinition())
	if strings.Contains(file.Name(), "/private/") && strings.HasSuffix(file.Name(), "_p.h") {
		return true
	}
	return false
}

// TODO
func is_qt_inner_class(cursor clang.Cursor) bool {
	return false
}

func in_namespace(cursor clang.Cursor) bool {
	return false
}

// TODO auto detect
// param cursor should be a class cursor
func is_protected_dtor_class(cursor clang.Cursor) bool {
	protectedDtors := map[string]int{
		"QTextCodec": 1, "QAccessibleInterface": 1, "QTextBlockGroup": 1,
		"QTextObject": 1, "QAccessibleWidget": 1,
		"QWebEngineSettings": 1, "QWebEngineHistory": 1, "QWebEngineUrlRequestInfo": 1,
		"QCameraExposure": 1, "QCameraFocus": 1, "QCameraImageProcessing": 1,
	}
	_, ok := protectedDtors[cursor.Spelling()]
	ok2 := false
	cursor.Visit(func(c, p clang.Cursor) clang.ChildVisitResult {
		switch c.Kind() {
		case clang.Cursor_Destructor:
			if c.AccessSpecifier() == clang.AccessSpecifier_Protected {
				ok2 = true
				return clang.ChildVisit_Break
			}
		}
		return clang.ChildVisit_Recurse
	})
	if ok2 != ok {
		log.Println("differenct detect result: ", ok, ok2, cursor.Spelling())
	}
	return ok
}

func is_qt_global_func(cursor clang.Cursor) bool {
	// qputenv,qsrand,qCompress
	reg := regexp.MustCompile(`q[A-Z].+`) // 需要生成的全局函数名正则规范
	reg = regexp.MustCompile(`q.+`)       // 需要生成的全局函数名正则规范
	// and is stdglobal scope?
	pcs := cursor.SemanticParent()
	if pcs.Kind() == clang.Cursor_Namespace &&
		(pcs.Spelling() == "Qt" || pcs.Spelling() == "QtPrivate" || pcs.Spelling() == "QtWebEngine" ||
			pcs.Spelling() == "QtAndroid" || pcs.Spelling() == "QtWin" || pcs.Spelling() == "QtMac") {
		return true
	}
	if cursor.Kind() == clang.Cursor_FunctionDecl && cursor.IsCursorDefinition() && cursor.IsFunctionInlined() {
		return true
	}
	if cursor.Kind() == clang.Cursor_FunctionDecl && !cursor.IsCursorDefinition() && !cursor.IsFunctionInlined() {
		return true
	}

	return reg.MatchString(cursor.Spelling())
}

func is_qstring_cls(retPlace string) bool {
	if funk.ContainsString(strings.FieldsFunc(retPlace, func(c rune) bool {
		return strings.Contains(" .*/", string(c))
	}), "QString") {
		return true
	}
	return false
}

func FuncHasLongDoubleArg(c clang.Cursor) bool {
	for i := int32(0); i < c.NumArguments(); i++ {
		if c.Argument(uint32(i)).Type().Kind() == clang.Type_LongDouble {
			return true
		}
	}
	return false
}

func TypeIsCharPtrPtr(ty clang.Type) bool {
	return (isPrimitivePPType(ty) && ty.PointeeType().PointeeType().Kind() == clang.Type_Char_S) ||
		(ty.Kind() == clang.Type_IncompleteArray && TypeIsCharPtr(ty.ElementType()))
}
func TypeIsCharPtr(ty clang.Type) bool {
	return ty.Kind() == clang.Type_Pointer && ty.PointeeType().Kind() == clang.Type_Char_S
}
func TypeIsUCharPtr(ty clang.Type) bool {
	return ty.Kind() == clang.Type_Pointer && ty.PointeeType().Kind() == clang.Type_UChar
}
func TypeIsBoolPtr(ty clang.Type) bool {
	return ty.Kind() == clang.Type_Pointer && ty.PointeeType().Kind() == clang.Type_Bool
}
func TypeIsVoidPtr(ty clang.Type) bool {
	return ty.Kind() == clang.Type_Pointer && ty.PointeeType().Kind() == clang.Type_Void
}
func TypeIsPtr(ty clang.Type) bool { return ty.Kind() == clang.Type_Pointer }
func TypeIsIntPtr(ty clang.Type) bool {
	return ty.Kind() == clang.Type_Pointer && ty.PointeeType().Kind() == clang.Type_Int
}

func TypeIsQFlags(ty clang.Type) bool {
	if ty.Kind() == clang.Type_Typedef &&
		strings.HasPrefix(ty.CanonicalType().Spelling(), "QFlags") &&
		strings.ContainsAny(ty.CanonicalType().Spelling(), "<>") {
		if strings.Contains(ty.Spelling(), "::") { // for QFlags<xxx> Qxxx::xxx
			return true
		}
	}
	return false
}

func TypeIsFuncPointer(ty clang.Type) bool {
	return strings.Contains(ty.Spelling(), " (*)(")
}

func TypeIsConsted(ty clang.Type) bool {
	return ty.IsConstQualifiedType() || strings.HasPrefix(ty.Spelling(), "const ")
}

var privClasses = map[string]int{"QV8Engine": 1, "QQmlComponentAttached": 1,
	"QQmlImageProviderBase": 1}

func rewriteOperatorMethodName(name string) string {
	if !strings.HasPrefix(name, "operator") {
		return name
	}
	replaces := []string{
		"&=", "_and_equal",
		"^=", "_caret_equal",
		"|=", "_or_equal",
		"+=", "_add_equal",
		"-=", "_minus_equal",
		"==", "_equal_equal",
		"!=", "_not_equal",
		"<<", "_left_shift",
		">>", "_right_shift",
		"[]", "_get_index",
		"()", "_fncall",
		"->", "_minus_greater",
		"<", "_less_than", ">", "_greater_than",
		"!", "_not", "=", "_equal",
		"&", "_and", "^", "_caret", "|", "_or", "~", "_around",
		"/", "_div", "*", "_mul", "-", "_minus", "+", "_add",
		" ", "_"}
	valiname := name
	for i := 0; i < len(replaces)/2; i += 1 {
		valiname = strings.Replace(valiname, replaces[i*2], replaces[i*2+1], -1)
	}
	return valiname
}

var fileCache = map[string]*os.File{}

func readSourceRange(sr clang.SourceRange) (string, string) {
	bfp, blineno, bcol, boffset := sr.Start().ExpansionLocation()
	efp, elineno, ecol, eoffset := sr.End().ExpansionLocation()
	log.Println(bfp.Name(), efp.Name(), len(bfp.Name()), len(efp.Name()), blineno, bcol, boffset, elineno, ecol, eoffset)

	if bfp.Name() == "" {
		return "", ""
	}
	if bfp.Name() != efp.Name() {
		log.Fatalln("wtf", bfp.Name())
	}

	var fph *os.File
	if fph_, ok := fileCache[bfp.Name()]; ok {
		fph = fph_
	} else {
		fph_, err := os.Open(bfp.Name())
		gopp.ErrPrint(err, bfp.Name())
		fph = fph_
		fileCache[bfp.Name()] = fph
	}
	if fph == nil {
		log.Fatalln("wtf", bfp.Name())
	}

	fph.Seek(int64(boffset), os.SEEK_SET)
	var buf = make([]byte, eoffset-boffset)
	n, err := fph.ReadAt(buf, int64(boffset))
	gopp.ErrPrint(err, n, boffset, eoffset)
	spelling := string(buf[:n])
	pos := strings.Index(spelling, " ")
	macroval := gopp.IfElseStr(pos > 0, spelling[pos+1:], "")
	if len(macroval) > 0 {
		macroty := ""
		if strings.HasPrefix(macroval, "\"") {
			macroty = "str"
		} else if strings.HasPrefix(macroval, "0x") {
			macroty = "num16"
		} else if gopp.IsNumberic(macroval) {
			macroty = "num10"
		}
		log.Println(bfp.Name(), n, spelling, len(macroval), macroval, macroty)
		if macroty != "" {
			return macroval, macroty
		}
	}

	return "", ""
}

func num_default_value(mth clang.Cursor) (n int) {
	for i := int32(0); i < mth.NumArguments(); i++ {
		argcs := mth.Argument(uint32(i))
		_, has := has_default_value(argcs)
		n += gopp.IfElseInt(has, 1, 0)
	}
	return
}

func has_default_value(arg clang.Cursor) (string, bool) {
	bfp, _, _, boffset := arg.Location().FileLocation()

	var fph *os.File
	if fph_, ok := fileCache[bfp.Name()]; ok {
		fph = fph_
	} else {
		fph_, err := os.Open(bfp.Name())
		gopp.ErrPrint(err, bfp.Name())
		fph = fph_
		fileCache[bfp.Name()] = fph
	}
	if fph == nil {
		log.Fatalln("wtf", bfp.Name())
	}

	fph.Seek(int64(boffset), os.SEEK_SET)
	s := ""
	hasdv := false
	leftb := 0
	buf := make([]byte, 1)
	for {
		_, err := fph.Read(buf)
		if err != nil {
			gopp.ErrPrint(err)
			break
		} else {
			if buf[0] == '(' {
				leftb += 1
			}
			if buf[0] == ',' || (buf[0] == ')' && leftb == 0) || buf[0] == '#' {
				break
			} else {
				if buf[0] == '=' {
					hasdv = true
				} else if hasdv {
					s += string(buf)
				}
			}
			if buf[0] == ')' {
				leftb -= 1
			}
		}
	}
	s = strings.TrimSpace(s)
	log.Println(arg.DisplayName(), arg.SemanticParent().Spelling(), hasdv, s)
	return s, hasdv
}

func getOrOpenCursorFile(c clang.Cursor) *os.File {
	bfp, _, _, _ := c.Location().FileLocation()

	var fph *os.File
	if fph_, ok := fileCache[bfp.Name()]; ok {
		fph = fph_
	} else {
		fph_, err := os.Open(bfp.Name())
		gopp.ErrPrint(err, bfp.Name())
		fph = fph_
		fileCache[bfp.Name()] = fph
	}
	if fph == nil {
		log.Fatalln("wtf", bfp.Name(), c.Spelling())
	}
	return fph
}

// support class/method/function cursor
// 从c的起始点往前推，找到第一段comment，当然不超过另一个c的定义（怎么判断？分号吗？可能有多个终止符号）
// 突然发现安装后的头文件中并没有注释，注释在cpp文件中的。
func readComment(c clang.Cursor) string {
	bfp, _, _, boffset := c.Location().FileLocation()

	var fph *os.File
	if fph_, ok := fileCache[bfp.Name()]; ok {
		fph = fph_
	} else {
		fph_, err := os.Open(bfp.Name())
		gopp.ErrPrint(err, bfp.Name())
		fph = fph_
		fileCache[bfp.Name()] = fph
	}
	if fph == nil {
		log.Fatalln("wtf", bfp.Name())
	}

	comment := ""
	buf := make([]byte, 0)
	ch := make([]byte, 1)
	log.Println("looping read comment:", boffset, c.Spelling(), bfp.Name())
	for off := boffset - 1; off >= 0; off-- {
		fph.ReadAt(ch, int64(off))
		if ch[0] == ';' {
			break
		}
		buf = append([]byte{ch[0]}, buf...)
		if len(buf) >= 2 {
			if buf[0] == '/' && buf[1] == '*' {
				comment = string(buf)
				log.Println("found comment:", c.Spelling(), len(comment))
				break
			}
		}
	}
	return comment
}

// since format: x.y
func sinceVer2Hex(since string) string {
	sepch := gopp.IfElseStr(strings.Contains(since, ","), ",", ".")
	parts := strings.Split(since, sepch)
	src := []byte{byte(gopp.MustInt(parts[0])), byte(gopp.MustInt(parts[1]))}
	hv := hex.EncodeToString(src)
	return fmt.Sprintf("0x%s00", hv)
}

var goqdocs = map[string]*goquery.Document{}

func refmtComment(s string, cmch string) string {
	return strings.Join(funk.Map(strings.Split(s, "\n"), func(s string) string { return cmch + s }).([]string), "\n")
}

// 查询方法/函数的注释，从doc/.html中
func queryComment(c clang.Cursor, qtdir, qtver string) string {
	mod := get_decl_mod(c)
	bfp, _, _, _ := c.Location().FileLocation()
	parts := strings.Split(bfp.Name(), "/")
	htmlFile := fmt.Sprintf("%s/Docs/Qt-%s/qt%s/%stml", qtdir, qtver, mod, parts[len(parts)-1])
	log.Println(bfp.Name(), "=>", htmlFile)

	switch c.Kind() {
	case clang.Cursor_CXXMethod, clang.Cursor_Constructor, clang.Cursor_Destructor,
		clang.Cursor_FunctionDecl:
		sltor := fmt.Sprintf("h3#%s.fn", c.Spelling())
		return queryCommentFromFile(htmlFile, c.Spelling(), sltor)
	case clang.Cursor_EnumDecl:
		sltor := fmt.Sprintf("h3#%s-enum.fn", c.Spelling())
		comment := queryCommentFromFile(htmlFile, c.Spelling(), sltor)
		if comment == "" {
			htmlFile = fmt.Sprintf("%s/Docs/Qt-%s/qtcore/qt.html", qtdir, qtver)
			comment = queryCommentFromFile(htmlFile, c.Spelling(), sltor)
		}
		return comment
	case clang.Cursor_ClassDecl, clang.Cursor_StructDecl:
		sltor := fmt.Sprintf("div#details")
		return queryCommentFromFile(htmlFile, c.Spelling(), sltor)
	}
	return "nonono"
}

func queryCommentFromFile(htmlFile string, name string, sltor string) string {
	if !gopp.FileExist(htmlFile) {
		return ""
	}

	var doco *goquery.Document
	if doco_, ok := goqdocs[htmlFile]; ok {
		doco = doco_
	} else {
		fp, err := os.Open(htmlFile)
		gopp.ErrPrint(err, htmlFile)
		doc, err := goquery.NewDocumentFromReader(fp)
		gopp.ErrPrint(err)
		fp.Close()
		doco = doc
	}

	slts := doco.Find(sltor)
	// log.Println(slts.Length(), sltor)
	if slts.Length() > 0 {
		n := slts.First().Nodes[0]
		comment := ""
		for n != nil {
			nn := n.NextSibling
			// log.Printf("%s, %+v\n", nn.Data, nn)
			// log.Println(goquery.NewDocumentFromNode(nn).Text())
			// log.Println(name, len(nn.Data), nn.Data)
			isdivtab := false
			for _, attr := range nn.Attr {
				if attr.Key == "class" && attr.Val == "table" {
					isdivtab = true
				}
			}

			if nn.Data == "p" || nn.Data == "pre" || nn.Data == "ul" {
				comment += "\n\n" + goquery.NewDocumentFromNode(nn).Text()
			} else if nn.Data == "div" && isdivtab {
				// omit table text
				comment += "\n\n" + goquery.NewDocumentFromNode(nn).Text()
			} else if strings.TrimSpace(nn.Data) != "" {
				break
			}

			n = nn
		}
		log.Println(name, ":", len(comment), comment)
		comment = strings.TrimSpace(comment)
		comment = strings.Replace(comment, "/*", "/-*", -1)
		comment = strings.Replace(comment, "*/", "*-/", -1)
		return comment
	}

	return ""
}

func extractEnumElem(comment string) (pureComment string, elems map[string]string) {
	elems = map[string]string{}
	lines := strings.Split(comment, "\n")
	exp := `(.+)::(.+)([0-9]+)(.+)`
	reg := regexp.MustCompile(exp)
	for _, line := range lines {
		if line == "ConstantValueDescription" {
			continue
		}
		if reg.MatchString(line) {
			mats := reg.FindAllStringSubmatch(line, -1)
			elems[mats[0][2]] = mats[0][4]
			continue
		}
		pureComment += line + "\n"
	}
	return
}

func getIncludeNameByModule(mod string) string {
	for name, _ := range modDepsAll {
		if strings.ToLower(mod) == strings.ToLower(name) &&
			strings.ToLower(mod) != name {
			return name
		}
	}
	return mod
}

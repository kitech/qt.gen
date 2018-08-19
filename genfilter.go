package main

import (
	"flag"
	"log"
	"regexp"
	"strings"

	"github.com/go-clang/v3.9/clang"
	funk "github.com/thoas/go-funk"
)

// any comment on this???
// see blow init func
var specifyClass string

func init() {
	flag.StringVar(&specifyClass, "gclass", specifyClass, "specify need generate one class")
}

type GenFilter interface {
	skipClass(cursor, parent clang.Cursor) bool
	skipMethod(cursor, parent clang.Cursor) bool
	skipArg(cursor, parent clang.Cursor) bool
	skipFunc(cursor clang.Cursor) bool
}

type GenFilterBase struct {
}

func (this *GenFilterBase) skipClass(cursor, parent clang.Cursor) bool {
	skip := this.skipClassImpl(cursor, parent)
	if strings.Contains(cursor.Spelling(), "QWidgetList") {
		// log.Fatalln(cursor.Spelling())
	}
	if strings.Contains(cursor.Spelling(), "QWidgetList") && skip > 0 {
		// log.Fatalln("skipped class:", skip)
	}
	return skip > 0
}

// TODO  拆分成多个小的过滤函数
func (this *GenFilterBase) skipClassImpl(cursor, parent clang.Cursor) int {
	cname := cursor.Spelling()
	prefixes := []string{
		"QMetaTypeId", "QTypeInfo", "QOpenGLFunctions",
		"QOpenGLExtraFunctions", "QOpenGLVersion", "QOpenGL",
		"QAbstract-", "QPrivate",
	}
	equals := []string{
		"QAbstractOpenGLFunctionsPrivate",
		"QOpenGLFunctionsPrivate",
		"QOpenGLExtraFunctionsPrivate",
		"QAnimationGroup-",
		"QMetaType",
		"QAtomicOpsSupport",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(cname, prefix) {
			return 1
		}
	}
	for _, equal := range equals {
		if equal == cname {
			return 2
		}
	}

	// 这个也许是因为qt有bug，也许是因为arch上的qt包有问题。QT_OPENGL_ES_2相关。
	if strings.HasPrefix(cname, "QOpenGLFunctions_") &&
		strings.Contains(cname, "CoreBackend") {
		return 3
	}
	if strings.HasPrefix(cname, "QOpenGLFunctions_") &&
		strings.Contains(cname, "DeprecatedBackend") {
		return 4
	}

	if !cursor.IsCursorDefinition() {
		// log.Println("filtered by not definition", cursor.Spelling())
		return 5
	}
	// pure virtual class check
	pure_virtual_class := is_pure_virtual_class(cursor)
	if pure_virtual_class {
		// return true
	}

	// if cursor.kind == clidx.CursorKind.CLASS_TEMPLATE: return True
	if cursor.SpecializedCursorTemplate().IsNull() == false {
		return 6
	}
	// inner class
	if cursor.SemanticParent().Kind() == clang.Cursor_StructDecl {
		return 7
	}
	if cursor.SemanticParent().Kind() == clang.Cursor_StructDecl {
		return 8
	}
	// test
	fixclasses := []string{"QDebug", "QNoDebug", "QDebugStateSaver", // "QFileDevice",
		"QLibraryInfo", "QInternal", "QAccessibleObject", "QAccessibleActionInterface",
		"QGraphicsObject"}
	for _, c := range fixclasses {
		if cursor.Spelling() == c {
			return 9
		}
	}
	if cname != "QString" {
		// return true
	}
	if cname != "QStringRef" {
		// return true
	}
	if cname != "QSysInfo" {
		// return true
	}
	if cname != "QCoreApplication" {
		// return true
	}
	if len(specifyClass) > 0 && cname != specifyClass {
		return 10
	}

	return 0
}

func (this *GenFilterBase) skipMethod(cursor, parent clang.Cursor) bool {
	skip := this.skipMethodImpl(cursor, parent)
	if cursor.Spelling() == "QApplication" {
	}
	log.Println(skip, cursor.Spelling(), parent.Spelling(), cursor.DisplayName(), skip)
	if skip > 0 {
		log.Println(skip, cursor.Spelling(), parent.Spelling(), cursor.DisplayName(), skip, cursor.AccessSpecifier())
		// os.Exit(0)
	}
	return skip > 0
}

func (this *GenFilterBase) skipMethodImpl(cursor, parent clang.Cursor) int {
	if cursor.AccessSpecifier() != clang.AccessSpecifier_Public {
		if cursor.AccessSpecifier() != clang.AccessSpecifier_Protected {
			if cursor.Spelling() == "QPaintDevice" {
				return 0
			}
			return 1
		}
	}

	cname := cursor.Spelling()
	metamths := []string{"qt_metacall", "qt_metacast", "qt_check_for_"}
	for _, mm := range metamths {
		if strings.HasPrefix(cname, mm) {
			return 2
		}
	}

	if funk.ContainsString([]string{"tr", "trUtf8", "data_ptr", "d_func"}, cname) {
		return 3
	}

	if strings.HasPrefix(cname, "operator") {
		// return 4
	}

	if funk.ContainsString([]string{"rend", "append", "insert", "rbegin", "prepend", "crend", "crbegin"}, cname) {
		return 5
	}

	// TODO return template container
	if funk.ContainsString([]string{"rawHeaderPairs", "rawHeaders"}, cname) {
		return 55
	}

	if cursor.IsVariadic() {
		return 6
	}
	// TODO move ctor and copy ctor?
	if cursor.CXXConstructor_IsCopyConstructor() {
		return 7
	}
	if cursor.CXXConstructor_IsMoveConstructor() {
		return 8
	}
	if strings.HasPrefix(cname, "operator") {
		if (parent.Spelling() == "QSignalBlocker" && cursor.Spelling() == "operator=") ||
			(parent.Spelling() == "QLoggingCategory" && cursor.Spelling() == "operator()") ||
			(parent.Spelling() == "QSemaphoreReleaser" && cursor.Spelling() == "operator=") {
			return 9
		}
	}

	//
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		if this.skipArg(cursor.Argument(uint32(idx)), cursor) {
			return 10
		}
	}

	if this.skipReturn(cursor.ResultType(), cursor) {
		return 11
	}

	return 0
}

func (this *GenFilterBase) skipFunc(cursor clang.Cursor) bool {
	if cursor.IsVariadic() {
		return true
	}
	if strings.Contains(cursor.Spelling(), "printf") {
		return true
	}
	if strings.Contains(cursor.DisplayName(), "QDebug") {
		return true
	}
	if strings.Contains(cursor.Spelling(), "qt_builtin_") {
		return true
	}
	if strings.Contains(cursor.Spelling(), "qustrlen") {
		return true
	}
	if strings.Contains(cursor.Spelling(), "_destructor") {
		return true
	}
	if this.skipReturn(cursor.ResultType(), cursor) {
		return true
	}
	return false
}

func (this *GenFilterBase) skipArg(cursor, parent clang.Cursor) bool {
	skip := this.skipArgImpl(cursor, parent)
	if skip > 0 {
		log.Println(skip, cursor.Type().Spelling(), cursor.Type().Kind().String(), cursor.Spelling(), parent.DisplayName())
	}
	return skip > 0
}

func (this *GenFilterBase) skipArgImpl(cursor, parent clang.Cursor) int {
	// C_ZN16QCoreApplication11aboutToQuitENS_14QPrivateSignalE(void *this_, QCoreApplication::QPrivateSignal a0)
	argTyBare := get_bare_type(cursor.Type())
	// log.Println(cursor.DisplayName(), cursor.DisplayName(), cursor.Type().Spelling(), is_qt_private_class(argTyBare.Declaration()), argTyBare.Spelling())
	if strings.Contains(argTyBare.Spelling(), "QPrivate") {
		return 1
	}
	if strings.HasSuffix(argTyBare.Spelling(), "Function") {
		if argTyBare.Spelling() == "QImageCleanupFunction" {
		} else {
			return 2
		}
	}
	if strings.HasSuffix(argTyBare.Spelling(), "Func") {
		return 3
	}
	if strings.HasSuffix(argTyBare.Spelling(), "Private") {
		return 4
	}
	if strings.HasSuffix(argTyBare.Spelling(), "DataPtr") {
		return 5
	}
	if _, ok := skipClasses[argTyBare.Spelling()]; !ok {
		if argTyBare.Kind() != clang.Type_Invalid && !isPrimitiveType(argTyBare) {
			// like Qt::WindowFlags form
			reg := regexp.MustCompile(`^Q.+::.*Flags$`)
			if reg.MatchString(argTyBare.Spelling()) {
			} else if argTyBare.Kind() == clang.Type_Typedef && strings.HasPrefix(argTyBare.CanonicalType().Spelling(), "QFlags<") {
			} else if argTyBare.Spelling() == "QImageCleanupFunction" {
			} else {
				log.Println(argTyBare.Spelling(), argTyBare.Kind().String(), argTyBare.CanonicalType().Spelling(), "skiped by skiped class")
				return 6
			}
		}
	}

	inenums := []string{
		"ComponentFormattingOptions",
		"FormattingOptions",
		"CategoryFilter",
		"KeyValues",
		"InterfaceFactory",
		"RootObjectHandler",
		"UpdateHandler",
		"QtMetaTypePrivate",
		"va_list",
	}
	for _, inenum := range inenums {
		if strings.Contains(cursor.Type().Spelling(), inenum) {
			return 7
		}
	}
	if cursor.Type().Spelling() == "Id" {
		return 8
	}
	// C_ZN18QThreadStorageDataC1EPFvPvE(void (*)(void *) func) {
	if cursor.Type().Spelling() == "void (*)(void *)" {
		return 9
	}

	if this.skipType(cursor.Type(), cursor) {
		return 10
	}

	return 0
}

func (this *GenFilterBase) skipReturn(ty clang.Type, cursor clang.Cursor) bool {
	skip := this.skipReturnImpl(ty, cursor)
	if skip > 0 {
		log.Println(skip, ty.Spelling(), cursor.DisplayName())
	}
	return skip > 0
}

func (this *GenFilterBase) skipReturnImpl(ty clang.Type, cursor clang.Cursor) int {
	log.Println(ty.Spelling(), cursor.Spelling(), ty.CanonicalType().Spelling())
	skips := []string{"QTimeZone::OffsetDataList", "QVariantAnimation::KeyValues",
		"QDebug", "QNoDebug", "QXmlStreamNamespaceDeclarations", "QXmlStreamNotationDeclarations",
		"QXmlStreamEntityDeclarations", "QGradientStops", "QOpenGLContext",
		"QAccessibleActionInterface", "QPlatformBackingStore", "QPlatformNativeInterface",
		"QPlatformOffscreenSurface", "QMatrix3x3", "QPagedPaintDevicePrivate",
		"QPlatformPixmap", "QPlatformScreen", "QPlatformSurface", "QTextDocumentPrivate",
		"QTextEngine", "QPlatformWindow", "QVulkanInstance", "QGraphicsEffectSource",
		"QGraphicsObject", "QPlatformMenu", "QPlatformMenuBar",
		"QOpenGLFramebufferObject", "QOpenGLShaderProgram",
		"QQmlWebChannelAttached"}
	for _, tn := range skips {
		log.Println(ty.Spelling(), cursor.Spelling(), ty.CanonicalType().Spelling(), ty.PointeeType().CanonicalType().Spelling(), ty.PointeeType().Declaration().Type().Spelling())
		if ty.Spelling() == tn || ty.PointeeType().Spelling() == tn ||
			ty.PointeeType().CanonicalType().Spelling() == tn ||
			ty.PointeeType().Declaration().Type().Spelling() == tn {
			return 1
		}
	}

	barety := get_bare_type(ty)
	if this.skipClass(barety.Declaration(), barety.Declaration().SemanticParent()) {
		// return 4
	}

	bareSpell := strings.Replace(ty.Spelling(), "const", "", -1)
	bareSpell = strings.Replace(bareSpell, "&", "", -1)
	bareSpell = strings.TrimSpace(bareSpell)

	isQListx := strings.HasPrefix(bareSpell, "QList<") &&
		funk.ContainsString([]string{"QUrl", "QSize", "QCameraInfo",
			"QObject *", "QQuickItem *", "QGraphicsItem *"},
			strings.TrimRight(bareSpell[6:], ">"))
	if strings.HasPrefix(bareSpell, "Q") {
		if strings.HasSuffix(bareSpell, "Map") ||
			// strings.HasSuffix(bareSpell, "List") ||
			// strings.HasSuffix(bareSpell, "Set") ||
			strings.HasSuffix(bareSpell, "Hash") {
			return 2
		}
		if strings.HasPrefix(bareSpell, "QList<") {
			if isQListx {
			} else {
				return 5
			}
		}
	}

	switch ty.Kind() {
	case clang.Type_Unexposed:
		if isQListx {
		} else {
			return 3
		}
	//case clang.Type_Pointer:
	case clang.Type_Void:
	default:
		if !isPrimitiveType(ty) {
			log.Println(ty.Kind().String(), ty.Spelling(), cursor.Spelling())
		}
	}
	return 0
}

func (this *GenFilterBase) skipType(ty clang.Type, cursor clang.Cursor) bool {
	skip := this.skipTypeImpl(ty, cursor)
	if skip > 0 {
		log.Println(skip, ty.Kind().Spelling(), ty.Spelling(), cursor.Spelling())
	}
	return skip > 0
}

func (this *GenFilterBase) skipTypeImpl(ty clang.Type, cursor clang.Cursor) int {

	switch ty.Kind() {
	case clang.Type_LValueReference:
		fallthrough
	case clang.Type_RValueReference:
		fallthrough
	case clang.Type_Pointer:
		// is template
		if ty.PointeeType().NumTemplateArguments() != -1 {
			return 1
		}
	case clang.Type_MemberPointer:
		return 2
	case clang.Type_Typedef:
		if false {
			log.Println(ty.Kind().Spelling(), ty.CanonicalType().Kind().Spelling())
		}
		return this.skipTypeImpl(ty.CanonicalType(), cursor)
	default:
		if ty.NumTemplateArguments() != -1 {
			// if strings.HasPrefix(ty.CanonicalType().Spelling(), "QFlags<") 这个会过滤掉太多方法
			if !strings.HasPrefix(ty.CanonicalType().Spelling(), "QFlags<") {
				return 3
			}
		}
	}

	return 0
}

/////
/*
base中过滤的项：
* 模板类声明
* 模板函数声明
* 模板方法声明
* QtPrivate::类，这个变动太大，要过滤掉
* 内部临时函数
*/
// 过滤原因的返回
type GenFilterBase2 struct{}

func (this *GenFilterBase2) skipClass(cursor, parent clang.Cursor) bool {
	skipn := this.skipClassImpl(cursor, parent)
	log.Println(skipn, cursor.Spelling(), parent.Spelling())
	return skipn > 0
}

// TODO  拆分成多个小的过滤函数
func (this *GenFilterBase2) skipClassImpl(cursor, parent clang.Cursor) int {
	cname := cursor.Spelling()
	prefixes := []string{
		"QMetaTypeId", "QTypeInfo", "QQmlTypeInfo", "QIntegerForSize",
		// "QOpenGLFunctions",
		// "QOpenGLExtraFunctions", "QOpenGLVersion", "QOpenGL",
		// "QAbstract-", "QPrivate",
	}
	equals := []string{
		"QAbstractOpenGLFunctionsPrivate",
		"QOpenGLFunctionsPrivate",
		"QOpenGLExtraFunctionsPrivate",
		"QOpenGLVersionFunctionsStorage",
		// "QAnimationGroup-",
		// "QMetaType",
		"QAtomicOpsSupport",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(cname, prefix) {
			return 1
		}
	}
	for _, equal := range equals {
		if equal == cname {
			return 2
		}
	}

	// 这个也许是因为qt有bug，也许是因为arch上的qt包有问题。QT_OPENGL_ES_2相关。
	if strings.HasPrefix(cname, "QOpenGLFunctions_") &&
		strings.Contains(cname, "CoreBackend") {
		// return 3
	}
	if strings.HasPrefix(cname, "QOpenGLFunctions_") &&
		strings.Contains(cname, "DeprecatedBackend") {
		// return 4
	}

	if !cursor.IsCursorDefinition() {
		// log.Println("filtered by not definition", cursor.Spelling())
		return 5
	}
	// pure virtual class check
	pure_virtual_class := is_pure_virtual_class(cursor)
	if pure_virtual_class {
		// return true
	}

	// if cursor.kind == clidx.CursorKind.CLASS_TEMPLATE: return True
	if cursor.SpecializedCursorTemplate().IsNull() == false {
		// return 6
	}
	// inner class
	if cursor.SemanticParent().Kind() == clang.Cursor_StructDecl {
		// return 7
	}
	if cursor.SemanticParent().Kind() == clang.Cursor_StructDecl {
		// return 8
	}
	// test
	fixclasses := []string{
		// "QDebug", "QNoDebug", "QDebugStateSaver", // "QFileDevice",
		// "QLibraryInfo", "QInternal", "QAccessibleObject", "QAccessibleActionInterface",
		// "QGraphicsObject",
	}
	for _, c := range fixclasses {
		if cursor.Spelling() == c {
			return 9
		}
	}
	if cname != "QString" {
		// return true
	}
	if cname != "QStringRef" {
		// return true
	}
	if cname != "QSysInfo" {
		// return true
	}
	if cname != "QCoreApplication" {
		// return true
	}
	if len(specifyClass) > 0 && cname != specifyClass {
		return 10
	}
	if parent.Spelling() == "QtPrivate" || parent.Spelling() == "QtMetaTypePrivate" {
		return 11
	}

	return 0
}

func (this *GenFilterBase2) skipMethod(cursor, parent clang.Cursor) bool {
	skipn := this.skipMethodImpl(cursor, parent)
	log.Println(skipn, cursor.Spelling(), parent.Spelling(), _cmgl.origin(cursor))
	return skipn > 0
}
func (this *GenFilterBase2) skipMethodImpl(cursor, parent clang.Cursor) int {
	if cursor.AccessSpecifier() == clang.AccessSpecifier_Invalid ||
		cursor.AccessSpecifier() == clang.AccessSpecifier_Private {
		return 1
	}
	if is_deleted_method(cursor, parent) {
		return 2
	}

	cname := cursor.Spelling()
	metamths := []string{
		// "qt_metacall", "qt_metacast",
		"qt_check_for_",
	}
	for _, mm := range metamths {
		if strings.HasPrefix(cname, mm) {
			return 2
		}
	}

	return 0
}

func (this *GenFilterBase2) skipArg(cursor, parent clang.Cursor) bool {
	return false
}
func (this *GenFilterBase2) skipFunc(cursor clang.Cursor) bool {
	return false
}

/////
type GenFilterInc struct {
	fltb *GenFilterBase2
}

func NewGenFilterInc() *GenFilterInc {
	this := &GenFilterInc{}
	this.fltb = &GenFilterBase2{}
	return this
}

func (this *GenFilterInc) skipClass(cursor, parent clang.Cursor) bool {
	bskip := this.fltb.skipClass(cursor, parent)
	return bskip
}

func (this *GenFilterInc) skipMethod(cursor, parent clang.Cursor) bool {
	bskip := this.fltb.skipMethod(cursor, parent)
	return bskip
}

func (this *GenFilterInc) skipArg(cursor, parent clang.Cursor) bool {
	bskip := this.fltb.skipArg(cursor, parent)
	return bskip
}
func (this *GenFilterInc) skipFunc(cursor clang.Cursor) bool {
	bskip := this.fltb.skipFunc(cursor)
	return bskip
}

/////
type GenFilterGo struct {
	GenFilterBase
}

func (this *GenFilterGo) skipMethod(cursor, parent clang.Cursor) bool {
	bskip := this.GenFilterBase.skipMethod(cursor, parent)
	return bskip
}

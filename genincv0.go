package main

import (
	"fmt"
	"gopp"
	"gopp/gods"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"unsafe"

	"github.com/go-clang/v3.9/clang"
	"github.com/therecipe/qt/internal/binding/parser"
	funk "github.com/thoas/go-funk"
	// "github.com/ianlancetaylor/demangle"
)

// revert to original pure inline generate, no wrapper

type GenerateInlinev0 struct {
	// TODO move to base
	filter   GenFilter
	tyconver *TypeConvertGo
	mangler  GenMangler

	methods   []clang.Cursor
	cp        *CodePager
	cpcs      *CodeFS // mod => file =>,
	argDesc   []string
	argtyDesc []string
	paramDesc []string

	GenBase
	hasMyCls bool
}

func NewGenerateInlinev0(qtdir, qtver string) *GenerateInlinev0 {
	this := &GenerateInlinev0{}
	this.qtdir, this.qtver = qtdir, qtver
	this.filter = &GenFilterInc{}
	this.mangler = NewIncMangler()
	this.tyconver = NewTypeConvertGo()

	this.GenBase.funcMangles = map[string]int{}

	this.cp = NewCodePager()
	this.initBlocks(this.cp)

	return this
}

func (this *GenerateInlinev0) initBlocks(cp *CodePager) {
	blocks := []string{"header", "main", "use", "ext", "body", "footer"}
	for _, block := range blocks {
		cp.AddPointer(block)
		cp.APf(block, "")
	}
}

func (this *GenerateInlinev0) genClass(cursor, parent clang.Cursor) {
	if false {
		log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
	}
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}
	clsctx := &GenClassContext{}
	clsctx.clscs = cursor
	clso, found := qdi.findClass(cursor.Spelling())
	if found {
		clsctx.clso = clso
	}

	this.isPureVirtualClass = false
	this.hasMyCls = false
	this.genFileHeader(clsctx, cursor, parent)
	this.walkClass(clsctx, cursor, parent)
	this.genProtectedCallbacks(clsctx, cursor, parent)
	this.genProxyClass(clsctx, cursor, parent)
	this.genMethods(clsctx, cursor, parent)
	this.final(cursor, parent)
}

func (this *GenerateInlinev0) final(cursor, parent clang.Cursor) {
	// log.Println(this.cp.ExportAll())
	this.saveCode(cursor, parent)

	this.cp = NewCodePager()
	this.initBlocks(this.cp)
}
func (this *GenerateInlinev0) saveCode(cursor, parent clang.Cursor) {
	// qtx{yyy}, only yyy
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}
	modname := strings.ToLower(filepath.Base(filepath.Dir(file.Name())))[2:]
	savefile := fmt.Sprintf("src/%s/%s.cxx", modname, strings.ToLower(cursor.Spelling()))

	ioutil.WriteFile(savefile, []byte(this.cp.ExportAll()), 0644)
}

func (this *GenerateInlinev0) saveCodeToFile(modname, file string) {
	// qtx{yyy}, only yyy
	savefile := fmt.Sprintf("src/%s/%s.cxx", modname, file)
	log.Println(savefile)

	// log.Println(this.cp.AllPoints())
	bcc := this.cp.ExportAll()
	if strings.HasPrefix(bcc, "//") {
		bcc = bcc[strings.Index(bcc, "\n"):]
	}
	ioutil.WriteFile(savefile, []byte(bcc), 0644)

}

func (this *GenerateInlinev0) genFileHeader(clsctx *GenClassContext, cursor, parent clang.Cursor) {
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}

	fullModname := filepath.Base(filepath.Dir(file.Name()))
	ftpath := strings.ToLower(fmt.Sprintf("%s/%s", fullModname, filepath.Base(file.Name())))
	if featname, ok := clts.qtreqcfgs[ftpath]; ok {
		this.cp.APf("header", "#ifndef QT_MINIMAL")
		this.cp.APf("header", "#include <%s/%sglobal.h>",
			fullModname, gopp.IfElseStr(fullModname == "QtCore", "q", strings.ToLower(fullModname)))
		this.cp.APf("header", "#if QT_CONFIG(%s)", featname)
		this.cp.APf("footer", "#endif // #if QT_CONFIG(%s)", featname)
		this.cp.APf("footer", "#endif // #ifndef QT_MINIMAL")
	}

	if clsctx.clso != nil && clsctx.clso.Since != "" {
		this.cp.APf("header", "// since %s", sinceVer2Hex(clsctx.clso.Since))
	}
	this.cp.APf("header", "// %s", fix_inc_name(file.Name()))
	this.cp.APf("header", "#ifndef protected")        // for combile source code, so with #ifdef
	this.cp.APf("header", "#define protected public") // for protected function call
	this.cp.APf("header", "#define private public")   // macos clang++ happy this
	this.cp.APf("header", "#endif")
	hbname := filepath.Base(file.Name())
	if strings.HasSuffix(hbname, "_impl.h") {
		this.cp.APf("header", "#include <%s.h>", hbname[:len(hbname)-7])
	} else {
		this.cp.APf("header", "#include <%s>", filepath.Base(file.Name()))
	}

	this.cp.APf("header", "#include <%s>", fullModname)
	this.cp.APf("header", "#include \"callback_inherit.h\"")
	this.cp.APf("header", "")
}

func (this *GenerateInlinev0) walkClass(clsctx *GenClassContext, cursor, parent clang.Cursor) {
	pureVirt := false
	methods := make([]clang.Cursor, 0)

	// pcursor := cursor
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_Constructor:
			pureVirt = pureVirt || cursor.CXXMethod_IsPureVirtual()
			fallthrough
		case clang.Cursor_Destructor:
			fallthrough
		case clang.Cursor_CXXMethod:
			pureVirt = pureVirt || cursor.CXXMethod_IsPureVirtual()
			if !this.filter.skipMethod(cursor, parent) {
				methods = append(methods, cursor)
			} else {
				// log.Println("filtered:", cursor.Spelling())
			}
		case clang.Cursor_UnexposedDecl:
			// log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
			file, line, col, _ := cursor.Location().FileLocation()
			if false {
				log.Println(file.Name(), line, col, file.Time())
			}
		default:
			if false {
				log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
			}
		}
		return clang.ChildVisit_Continue
	})

	if !pureVirt {
		pureVirt = is_pure_virtual_class(cursor)
	}
	isabstract := cursor.CXXRecord_IsAbstract()
	if isabstract != pureVirt {
		log.Println("abstract class detect not match", cursor.Spelling(), pureVirt, isabstract)
	}
	this.isPureVirtualClass = pureVirt
	this.cp.APf("header", "// %s is pure virtual: %v %v",
		cursor.Spelling(), pureVirt, isabstract)
	this.methods = methods
}

func (this *GenerateInlinev0) genProtectedCallbacks(clsctx *GenClassContext, cursor, parent clang.Cursor) {
	log.Println("process class:", len(this.methods), cursor.Spelling())
	for _, cursor := range this.methods {
		parent := cursor.SemanticParent()
		// log.Println(cursor.Kind().String(), cursor.DisplayName())

		if cursor.AccessSpecifier() == clang.AccessSpecifier_Protected {
			this.genProtectedCallback(cursor, parent)
		}
	}

	this.cp.APf("main", "")
}

// 非静态类
// 本类Ctor
// 本类protected  virtual方法
// 本类 纯虚方法
// 本类的所有父类 未实现的纯虚方法
// 当然要排重，排重时注意const限定符也是检测唯一签名相关的
func (this *GenerateInlinev0) collectProxyMethod(clsctx *GenClassContext, cursor, parent clang.Cursor) (pxymths []clang.Cursor) {
	mthmap := gods.NewHashMap()   // name+constflag => clang.Cursor
	mthlst := gods.NewArrayList() // name+constflag
	mthkey := func(c clang.Cursor) string {
		k := c.Spelling() + fmt.Sprintf("-c%v", c.CXXMethod_IsConst())
		for i := int32(0); i < c.NumArguments(); i++ {
			k += "-" + c.Argument(uint32(i)).Type().Spelling()
		}
		return k
	}
	deletemth := func(c clang.Cursor) {
		key1 := mthkey(c)
		idx, ival := mthlst.Find(func(index int, value interface{}) bool {
			key2 := mthkey(value.(clang.Cursor))
			return key2 == key1
		})
		if ival != nil {
			mthlst.Remove(idx)
		} else {
			log.Println("delete failed:", idx, key1)
		}
	}

	bcs := find_base_classes(cursor)
	bcs2 := append(bcs, cursor)
	// 本类与基类的纯虚方法
	for i := len(bcs2) - 1; i >= 0; i-- {
		bci := bcs2[i]
		bci.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
			switch cursor.Kind() {
			case clang.Cursor_CXXMethod:
				if cursor.CXXMethod_IsStatic() {
					break
				}
				key := mthkey(cursor)
				if cursor.CXXMethod_IsPureVirtual() {
					if c, ok := mthmap.Get(key); ok {
						mthmap.Remove(key)
						deletemth(c.(clang.Cursor))
					}
					mthmap.Put(key, cursor)
					mthlst.Add(cursor)
				}
			}
			return clang.ChildVisit_Continue
		})
	}

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_CXXMethod:
			if cursor.CXXMethod_IsStatic() {
				break
			}
			key := mthkey(cursor)
			if !cursor.CXXMethod_IsPureVirtual() {
				if c, ok := mthmap.Get(key); ok {
					mthmap.Remove(key)
					deletemth(c.(clang.Cursor))
				}
			}
			if cursor.CXXMethod_IsVirtual() && cursor.AccessSpecifier() == clang.AccessSpecifier_Protected {
				if c, ok := mthmap.Get(key); ok {
					mthmap.Remove(key)
					deletemth(c.(clang.Cursor))
				}
				mthmap.Put(key, cursor)
				mthlst.Add(cursor)
			}

		case clang.Cursor_Constructor:
			if cursor.AccessSpecifier() == clang.AccessSpecifier_Public && !is_deleted_method(cursor, parent) {
				mthlst.Add(cursor)
			}
		case clang.Cursor_Destructor: // not need
		}
		return clang.ChildVisit_Continue
	})

	// gopp.Assert(mthmap.Size() == mthlst.Size(), "", mthmap.Size(), mthlst.Size(), cursor.Spelling())
	mthlst.Each(func(index int, value interface{}) { pxymths = append(pxymths, value.(clang.Cursor)) })
	return
}

func (this *GenerateInlinev0) genProxyClass(clsctx *GenClassContext, cursor, parent clang.Cursor) {
	this.hasVirtualProtected = false
	if is_deleted_class(cursor) {
		return
	}
	isqobjcls := has_qobject_base_class(cursor)
	_ = isqobjcls
	// 需要proxy的类：QObject的子类

	this.hasMyCls = true
	this.cp.APf("main", "struct qt_meta_stringdata_My%s_t {", cursor.Spelling())
	this.cp.APf("main", "  QByteArrayData data[1];")
	this.cp.APf("main", "  char stringdata0[%d];", 2+len(cursor.Spelling())+1)
	this.cp.APf("main", "};")
	this.cp.APf("main", "#define QT_MOC_LITERAL(idx, ofs, len) \\")
	this.cp.APf("main", "  Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \\")
	this.cp.APf("main", "  qptrdiff(offsetof(qt_meta_stringdata_My%s_t, stringdata0) + ofs \\", cursor.Spelling())
	this.cp.APf("main", "  - idx * sizeof(QByteArrayData)) \\")
	this.cp.APf("main", "  )")
	this.cp.APf("main", "static const qt_meta_stringdata_My%s_t qt_meta_stringdata_My%s = {", cursor.Spelling(), cursor.Spelling())
	this.cp.APf("main", "   {")
	this.cp.APf("main", "  QT_MOC_LITERAL(0, 0, %d), // \"My%s\"", 2+len(cursor.Spelling()), cursor.Spelling())
	this.cp.APf("main", "  },")
	this.cp.APf("main", "  \"My%s\"", cursor.Spelling())
	this.cp.APf("main", "};")
	this.cp.APf("main", "#undef QT_MOC_LITERAL")
	this.cp.APf("main", "static const uint qt_meta_data_My%s[] = {", cursor.Spelling())
	this.cp.APf("main", "  // content:")
	this.cp.APf("main", "  7,       // revision")
	this.cp.APf("main", "  0,       // classname")
	this.cp.APf("main", "  0,   0, // classinfo")
	this.cp.APf("main", "  0,   0, // methods")
	this.cp.APf("main", "  0,    0, // properties")
	this.cp.APf("main", "  0,    0, // enums/sets")
	this.cp.APf("main", "  0,    0, // constructors")
	this.cp.APf("main", "  0,       // flags")
	this.cp.APf("main", "  0,       // signalCount")
	this.cp.APf("main", "  0        // eod")
	this.cp.APf("main", "};")
	this.cp.APf("main", "class Q_DECL_EXPORT My%s : public %s {", cursor.Spelling(), cursor.Type().Spelling())
	if isqobjcls {
		this.cp.APf("main", "public: // Q_OBJECT")
		this.cp.APf("main", "/*static*/ QMetaObject staticMetaObject = {{&%s::staticMetaObject,", cursor.Spelling())
		this.cp.APf("main", "  qt_meta_stringdata_My%s.data,", cursor.Spelling())
		this.cp.APf("main", "  qt_meta_data_My%s,", cursor.Spelling())
		this.cp.APf("main", "  qt_static_metacall, nullptr, nullptr")
		this.cp.APf("main", "}};")
		this.cp.APf("main", "virtual const QMetaObject *metaObject() const override {")
		this.cp.APf("main", "  int handled = 0;")
		this.cp.APf("main", "  auto irv = callbackAllInherits_fnptr((void*)this, (char*)\"metaObject\", &handled, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0);")
		this.cp.APf("main", "   if (handled) { return (QMetaObject*)irv; }")
		this.cp.APf("main", "  return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject; ")
		this.cp.APf("main", "}")
		this.cp.APf("main", "virtual void *qt_metacast(const char *_clname) override {")
		this.cp.APf("main", "  int handled = 0;")
		this.cp.APf("main", "  auto irv = callbackAllInherits_fnptr((void*)this, (char*)\"qt_metacast\", &handled, 1, (uint64_t)_clname, 0, 0, 0, 0, 0, 0, 0, 0, 0);")
		this.cp.APf("main", "   if (handled) { return (void*)irv; }")
		this.cp.APf("main", "  if (!_clname) return nullptr;")
		this.cp.APf("main", "  if (!strcmp(_clname, qt_meta_stringdata_My%s.stringdata0))", cursor.Spelling())
		this.cp.APf("main", "      return static_cast<void*>(this);")
		this.cp.APf("main", "  return %s::qt_metacast(_clname);", cursor.Spelling())
		this.cp.APf("main", "}")
		this.cp.APf("main", "virtual int qt_metacall(QMetaObject::Call _c, int _id, void **_a) override {")
		this.cp.APf("main", "   _id = %s::qt_metacall(_c, _id, _a);", cursor.Spelling())
		this.cp.APf("main", "   if (_id < 0 ) return _id;")
		this.cp.APf("main", "   if (qt_metacall_fnptr != 0) {")
		this.cp.APf("main", "      return qt_metacall_fnptr(this, _c, _id, _a);")
		this.cp.APf("main", "   }")
		this.cp.APf("main", "   int handled = 0;")
		this.cp.APf("main", "   auto irv = callbackAllInherits_fnptr((void*)this, (char*)\"qt_metacall\", &handled, 3, (uint64_t)_c, (uint64_t)_id, (uint64_t)_a, 0, 0, 0, 0, 0, 0, 0);")
		this.cp.APf("main", "   if (handled) { return (int)irv; }")
		this.cp.APf("main", "   return _id;")
		this.cp.APf("main", "  }")
		this.cp.APf("main", "/*static*/ inline QString tr(const char *s, const char *c = nullptr, int n = -1)")
		this.cp.APf("main", "{ return staticMetaObject.tr(s, c, n); }")
		this.cp.APf("main", "/*static*/ inline QString trUtf8(const char *s, const char *c = nullptr, int n = -1)")
		this.cp.APf("main", " { return staticMetaObject.tr(s, c, n); }")
		this.cp.APf("main", "Q_DECL_HIDDEN_STATIC_METACALL static void qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a){")
		this.cp.APf("main", "  int handled = 0;")
		this.cp.APf("main", "  auto irv = callbackAllInherits_fnptr((void*)_o, (char*)\"qt_static_metacall\", &handled, 4, (uint64_t)_o, (uint64_t)_c, (uint64_t)_id, (uint64_t)_a, 0, 0, 0, 0, 0, 0);")
		this.cp.APf("main", "}")
		this.cp.APf("main", "private: struct QPrivateSignal {};")
		this.cp.APf("main", "")

		this.cp.APf("main", "public:")
		this.cp.APf("main", "  void* (*qt_metacast_fnptr)(void*, char*) = nullptr;")
		this.cp.APf("main", "  int (*qt_metacall_fnptr)(QObject *, QMetaObject::Call, int, void **) = nullptr;")
	}

	this.cp.APf("main", "public:")
	if !is_protected_dtor_class(cursor) {
		this.cp.APf("main", "  virtual ~My%s() {}", cursor.Spelling())
	}

	// override 目标，1. 让该类能够new, 2. 能够在binding端override可以override的方法
	overrideMethods := this.collectProxyMethod(clsctx, cursor, parent)
	proxyedMethods := []clang.Cursor{} // 这个要生成相应的公开调用封装函数
	for _, mcs := range overrideMethods {
		if mcs.Kind() == clang.Cursor_Constructor {
			this.cp.APf("main", "// %s %s", mcs.ResultType().Spelling(), mcs.DisplayName())
			this.genArgsPxy(mcs, cursor)
			argStr := strings.Join(this.argDesc, ", ")
			this.genParamsPxy(mcs, cursor)
			paramStr := strings.Join(this.paramDesc, ", ")
			if len(argStr) > 0 {
				// argStr = ", " + argStr
			}
			this.cp.APf("main", "My%s(%s) : %s(%s) {}", mcs.Spelling(),
				argStr, cursor.Type().Spelling(), paramStr)
			continue
		}
		if funk.Contains([]string{"metaObject", "qt_metacast", "qt_metacall"}, mcs.Spelling()) {
			continue
		}

		if mcs.AccessSpecifier() == clang.AccessSpecifier_Protected {
			this.hasVirtualProtected = true
		}
		rety := mcs.ResultType()
		this.cp.APf("main", "// %s", strings.Join(this.getFuncQulities(mcs), " "))
		this.cp.APf("main", "// [%d] %s %s", rety.SizeOf(), rety.Spelling(), mcs.DisplayName())

		// gen projected methods
		proxyedMethods = append(proxyedMethods, mcs)
		this.genArgsPxy(mcs, cursor)
		argStr := strings.Join(this.argDesc, ", ")
		this.genParamsPxy(mcs, cursor)
		paramStr := strings.Join(this.paramDesc, ", ")
		argStr2 := argStr
		argStr2 = gopp.IfElseStr(len(argStr2) > 0, ", "+argStr2, argStr2)
		paramStr2 := paramStr
		paramStr2 = gopp.IfElseStr(len(paramStr2) > 0, ", "+paramStr2, paramStr2)
		argtyStr := strings.Join(this.argtyDesc, ", ")
		argtyStr = gopp.IfElseStr(len(argtyStr) > 0, ", "+argtyStr, argtyStr)

		this.genArgs(mcs, cursor)
		argtyStr3 := strings.Join(this.argtyDesc, ", ")
		argtyStr3 = gopp.IfElseStr(len(argtyStr3) > 0, ", "+argtyStr3, argtyStr3)
		this.genParamsPxyCall(mcs, cursor)
		paramStr3 := strings.Join(this.paramDesc, ", ")
		paramStr3 = gopp.IfElseStr(len(paramStr3) > 0, ", "+paramStr3, paramStr3)

		prms := []string{}
		prms = append(prms, fmt.Sprintf("%d", mcs.NumArguments()))
		for i := int32(0); i < mcs.NumArguments(); i++ {
			arg := mcs.Argument(uint32(i))
			argty := arg.Type()
			// prmname := gopp.IfElseStr(arg.Spelling() == "", fmt.Sprintf("arg%d", i), arg.Spelling())
			prmname := this.genParamRefName(arg, mcs, int(i))
			switch argty.Kind() {
			case clang.Type_Record:
				prms = append(prms, "(uint64_t)&"+prmname)
			case clang.Type_LValueReference:
				prms = append(prms, "(uint64_t)&"+prmname)
			case clang.Type_Double, clang.Type_Float:
				prms = append(prms, "(uint64_t)&"+prmname)
			case clang.Type_Typedef:
				switch argty.Spelling() {
				case "qreal":
					prms = append(prms, "(uint64_t)&"+prmname)
				default:
					prms = append(prms, "(uint64_t)"+prmname)
				}
			case clang.Type_Unexposed: // QList<...
				prms = append(prms, "(uint64_t)&"+prmname)
			default:
				log.Println(argty.Kind(), arg.DisplayName(), arg.Spelling())
				prms = append(prms, "(uint64_t)"+prmname)
			}
		}
		for i := mcs.NumArguments(); i < 10; i++ {
			prms = append(prms, "0")
		}
		prmStr4 := strings.Join(prms, ", ")

		// TODO if anyway to know the binding peer if override a protected method, then can improve it
		constfix := gopp.IfElseStr(mcs.CXXMethod_IsConst(), "const", "")
		overrided := gopp.IfElseStr(mcs.CXXMethod_IsVirtual() || mcs.CXXMethod_IsPureVirtual(), "override", "")
		this.cp.APf("main", "  virtual %s %s(%s) %s %s {", rety.Spelling(), mcs.Spelling(), argStr, constfix, overrided)
		this.cp.APf("main", "    int handled = 0;")
		this.cp.APf("main", "    auto irv = callbackAllInherits_fnptr((void*)this, (char*)\"%s\", &handled, %s);",
			mcs.Spelling(), prmStr4)
		this.cp.APf("main", "    if (handled) {")
		switch rety.Kind() {
		case clang.Type_Void:
		case clang.Type_Record:
			this.cp.APf("main", "    if (irv == 0) { return (%s){};}", rety.Spelling())
			this.cp.APf("main", "    auto prv = (%s*)(irv); auto orv = *prv; delete(prv); return orv;", rety.Spelling())
		case clang.Type_Elaborated, clang.Type_Enum:
			if rety.CanonicalType().Kind() == clang.Type_Record &&
				rety.CanonicalType().SizeOf() == int64(unsafe.Sizeof(uint64(0))) {
				// for QAccessable::State like struct Elaborated type
				this.cp.APf("main", "    return *(%s*)(&irv);", rety.Spelling())
			} else {
				this.cp.APf("main", "    return (%s)(int)(irv);", rety.Spelling())
			}
		case clang.Type_Typedef:
			log.Println(rety.Spelling(), rety.CanonicalType().Kind(), rety.CanonicalType().Spelling(), rety.ClassType().Kind(), rety.ClassType().Spelling())
			if TypeIsQFlags(rety) {
				this.cp.APf("main", "    return (%s)(int)(irv);", rety.Spelling())
			} else if rety.CanonicalType().Kind() == clang.Type_Record {
				this.cp.APf("main", "    return *(%s*)(irv);", rety.Spelling())
			} else {
				this.cp.APf("main", "    return (%s)(irv);", rety.Spelling())
			}
		case clang.Type_Unexposed: // QList<...
			this.cp.APf("main", "    if (irv == 0) { return (%s){};}", rety.Spelling())
			this.cp.APf("main", "    auto prv = (%s*)(irv); auto orv = *prv; delete(prv); return orv;", rety.Spelling())
		default:
			this.cp.APf("main", "    return (%s)(irv);", rety.Spelling())
		}
		this.cp.APf("main", "      // %s %s %s", rety.Kind().String(), rety.CanonicalType().Kind().String(), rety.CanonicalType().Spelling())
		this.cp.APf("main", "    } else {")
		// TODO check return and convert return if needed
		ispurevirt := mcs.CXXMethod_IsPureVirtual()
		mthhasimpl := false
		if mcs.SemanticParent().Equal(cursor) && ispurevirt { // 当前类
		} else if !mcs.SemanticParent().Equal(cursor) && ispurevirt { // 父类
		} else {
			mthhasimpl = true
		}
		if mcs.ResultType().Kind() == clang.Type_Void {
			if mcs.Spelling() == "paintEvent" && cursor.Spelling() == "QWidget" {
				this.cp.APf("main", "    QStyleOption opt; opt.init(this); QPainter p(this);")
				this.cp.APf("main", "    style()->drawPrimitive(QStyle::PE_Widget, &opt, &p, this);")
			} else if mthhasimpl {
				this.cp.APf("main", "    %s::%s(%s);", cursor.Spelling(), mcs.Spelling(), paramStr)
			} else {
				this.cp.APf("main", "    // %s::%s(%s);", cursor.Spelling(), mcs.Spelling(), paramStr)
			}
		} else {
			if mthhasimpl {
				this.cp.APf("main", "    return %s::%s(%s);", cursor.Spelling(), mcs.Spelling(), paramStr)
			} else {
				if rety.Kind() == clang.Type_LValueReference {
					this.cp.APf("main", "    auto orv = (%s){}; return orv;", rety.PointeeType().Spelling())
				} else {
					this.cp.APf("main", "    return (%s){};", rety.Spelling())
				}
			}
		}
		this.cp.APf("main", "  }")
		this.cp.APf("main", "  }")
		this.cp.APf("main", "")
	}
	this.cp.APf("main", "};")
	this.cp.APf("main", "")

	if isqobjcls {
		this.cp.APf("main", "extern \"C\" Q_DECL_EXPORT")
		this.cp.APf("main", "void* C_%s_init_staticMetaObject(void* this_, void* strdat, void* dat, void* smcfn, void* mcastfn, void* mcallfn) {", cursor.Spelling())
		this.cp.APf("main", "  My%s* qo = (My%s*)(this_);", cursor.Spelling(), cursor.Spelling())
		this.cp.APf("main", "  QMetaObject* qmo = &qo->staticMetaObject;")
		this.cp.APf("main", "  qmo->d.stringdata = decltype(qmo->d.stringdata)(strdat);")
		this.cp.APf("main", "  qmo->d.data = decltype(qmo->d.data)(dat);")
		this.cp.APf("main", "  qmo->d.static_metacall = decltype(qmo->d.static_metacall)(smcfn);")
		this.cp.APf("main", "  qo->qt_metacast_fnptr = decltype(qo->qt_metacast_fnptr)(mcastfn);")
		this.cp.APf("main", "  qo->qt_metacall_fnptr = decltype(qo->qt_metacall_fnptr)( mcallfn);")
		this.cp.APf("main", "  return qmo;")
		this.cp.APf("main", "}")
		this.cp.APf("main", "")
	}

	// a hotfix
	if this.hasVirtualProtected && cursor.Spelling() == "QVariant" {
		this.hasVirtualProtected = false
	}
	this.genMethodsProxyed(clsctx, proxyedMethods)
}

func (this *GenerateInlinev0) genMethodsProxyed(clsctx *GenClassContext, methods []clang.Cursor) {
	for midx, method := range methods {
		this.genMethodProxyed(clsctx, method, midx)
	}
}

func (this *GenerateInlinev0) genMethodProxyed(clsctx *GenClassContext, cursor clang.Cursor, midx int) {

	this.genMethodHeader(clsctx, cursor, cursor.SemanticParent())
	// this.cp.APf("main", "")

	if !cursor.CXXMethod_IsPureVirtual() {
		parentSelector := ""
		parentSelector = fmt.Sprintf("%s::", cursor.SemanticParent().Spelling())
		this.genNonStaticMethod(clsctx, cursor, cursor.SemanticParent(), parentSelector)
	}
}

func (this *GenerateInlinev0) genMethods(clsctx *GenClassContext, cursor, parent clang.Cursor) {
	this.cp.APf("header", "// %s has virtual projected: %v", cursor.Spelling(), this.hasVirtualProtected)
	log.Println("process class:", len(this.methods), cursor.Spelling())

	cg := clcg // from temp global
	seeDtor := false
	for _, cursor := range this.methods {
		parent := cursor.SemanticParent()
		// log.Println(cursor.Kind().String(), cursor.DisplayName())
		if cursor.AccessSpecifier() == clang.AccessSpecifier_Protected {
			continue
		}

		isdeleted := cursor.IsDeleted()
		funcprototy := cursor.GetFunctionProtoType()
		log.Println(parent.Spelling(), cursor.Spelling(), isdeleted, funcprototy)

		hastpl := hasTmplArgRet(cursor)
		if !hastpl {
			fni := cg.ArrangeCXXMethodType(cursor, parent)
			// fni = cg.ArrangeFreeFunctionType(cursor)
			fnty := cg.ConvertFreeFunctionType(cursor)
			_ = fnty
			irargsz := cg.LLVMNumParams(cursor)
			issret := cursor.IsStructRet1(fni)
			argsz := cursor.ArgSize(fni)
			argsz1 := cursor.NumArguments()
			argsz2 := cursor.LLVMNumElements(fni)
			retriv := cursor.ResultType().IsTriviallyCopyableType()
			retriv2 := cursor.ResultType().IsTrivialType()
			retkd := cursor.ABIArgInfoKind(fni, -1)
			argkds := []clang.CGABIArgInfoKind{}
			argtys := []string{}
			for i := 0; i < argsz; i++ {
				argkds = append(argkds, cursor.ABIArgInfoKind(fni, i))
				argtys = append(argtys, cursor.ABIArgInfoType(fni, i).Spelling())
			}

			log.Println(parent.Spelling(), cursor.Spelling(), fni, "issret", issret, "aikind", retkd, argkds, cursor.ResultType().Spelling(), argsz, argsz1, "argstnum", argsz2, irargsz, cursor.Type().Spelling(), retriv, retriv2)
			if issret {
				//log.Panicln("got")
			}
		}

		if !cursor.IsFunctionInlined() {
			// 有些在类内没声明inline，但在类外又实现了inline的函数体, QWidget.resize
			defbody := cursor.Definition()
			if !defbody.IsFunctionInlined() {
				continue
			}
		}

		this.genMethodHeader(clsctx, cursor, parent)
		switch cursor.Kind() {
		case clang.Cursor_Constructor:
			this.genCtor(clsctx, cursor, parent)
		case clang.Cursor_Destructor:
			seeDtor = true
			this.genDtor(clsctx, cursor, parent)
		default:
			if cursor.CXXMethod_IsStatic() {
				this.genStaticMethod(clsctx, cursor, parent)
			} else {
				this.genNonStaticMethod(clsctx, cursor, parent, "")
			}
		}
	}
	if !seeDtor && !is_deleted_class(cursor) && !is_protected_dtor_class(cursor) {
		this.genDtorNotsee(cursor, parent)
	}
}

// TODO move to base
func (this *GenerateInlinev0) genMethodHeader(clsctx *GenClassContext, cursor, parent clang.Cursor) {
	qualities := this.getFuncQulities(cursor)
	if len(qualities) > 0 {
		this.cp.APf("main", "// %s", strings.Join(qualities, " "))
	}

	funco, found := (*parser.Function)(nil), false
	if clsctx.clso != nil {
		funco, found = qdi.findMethod(clsctx.clso, cursor)
		if found && funco.Since != "" {
			this.cp.APf("main", "// since %s", funco.Since)
		}
	}
	if cursor.Spelling() == "layoutChanged" && parent.Spelling() == "QAbstractItemModel" {
		// log.Fatalln(found, funco == nil)
	}

	file, lineno, _, _ := cursor.Location().FileLocation()
	this.cp.APf("main", "// %s:%d", fix_inc_name(file.Name()), lineno)
	consted := gopp.IfElseStr(cursor.CXXMethod_IsConst(), "const", "")
	this.cp.APf("main", "// [%d] %s %s %s",
		cursor.ResultType().SizeOf(), cursor.ResultType().Spelling(),
		cursor.DisplayName(), consted)
	this.cp.APf("main", "// (%d)%s (%d)%s",
		len(this.mangler.crc32p(cursor)), this.mangler.crc32p(cursor),
		len(this.mangler.origin(cursor)), this.mangler.origin(cursor))
}
func (this *GenerateInlinev0) genMethodFooter(clsctx *GenClassContext, cursor, parent clang.Cursor) {
	// this.cp.APf("main", "extern \"C\" Q_DECL_EXPORT")
	// this.cp.APf("main", "void* %s = (void*)&%s;", this.mangler.crc32p(cursor), this.mangler.convTo(cursor))
}
func (this *GenerateInlinev0) mthFnptrDesc(cursor, parent clang.Cursor) string {
	// mgname := this.mangler.origin(cursor)
	// like this: (int(QString::*)() const) &QString::count;
	filt := cursor.DisplayName() // no return type
	pos := strings.Index(filt, "(")
	signoret := filt[pos:]
	consted := gopp.IfElseStr(cursor.CXXMethod_IsConst(), "const", "")
	res := fmt.Sprintf("(%s (%s::*)%s %s) &%s::%s",
		cursor.ResultType().Spelling(), parent.Spelling(), signoret,
		consted, parent.Spelling(), cursor.Spelling())
	return res
}

func (this *GenerateInlinev0) genCtor(clsctx *GenClassContext, cursor, parent clang.Cursor) {
	// log.Println(this.mangler.convTo(cursor))
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	pparent := parent.SemanticParent()
	log.Println(cursor.Spelling(), parent.DisplayName(),
		cursor.SemanticParent().DisplayName(), cursor.LexicalParent().DisplayName(),
		pparent.Spelling(), parent.CanonicalCursor().DisplayName())

	funco, found := (*parser.Function)(nil), false
	if clsctx.clso != nil {
		funco, found = qdi.findMethod(clsctx.clso, cursor)
	}

	if found && funco.Since != "" {
		this.cp.APf("main", "#if QT_VERSION >= %s", sinceVer2Hex(funco.Since))
	}
	// this.cp.APf("main", "extern \"C\" Q_DECL_EXPORT")
	this.cp.APf("main", "void* %s(%s) {", this.mangler.crc32p(cursor), argStr)
	pxyclsp := ""
	if !is_deleted_class(parent) && this.hasVirtualProtected {
		this.cp.APf("main", "  auto _nilp = (My%s*)(0);", parent.Spelling())
		pxyclsp = "My"
		// pxyclsp = "" // TODO
	}
	isobjsub := has_qobject_base_class(parent)
	pureVirtRetstr := gopp.IfElseStr(this.isPureVirtualClass, "0; //", "")
	pureVirtRetstr = gopp.IfElseStr(this.isPureVirtualClass || !this.hasMyCls, "0; //", "")
	// TODO 要判断的，1,是否能加My前缀，2,不能加的情况，2-1,是否能new，2-2,不能new要加注册
	if isobjsub && !strings.ContainsAny(parent.Type().Spelling(), "<>") &&
		!funk.ContainsString([]string{"QAbstractEventDispatcher"}, parent.Type().Spelling()) {
		pxyclsp = "My"
		pureVirtRetstr = ""
	}

	if strings.HasPrefix(pparent.Spelling(), "Qt") {
		if pxyclsp == "" {
			this.cp.APf("main", "  return %s new %s::%s(%s);", pureVirtRetstr, pparent.Spelling(), parent.Spelling(), paramStr)
		} else {
			this.cp.APf("main", "  return %s new %s%s(%s);", pureVirtRetstr, pxyclsp, parent.Spelling(), paramStr)
		}
	} else {
		this.cp.APf("main", "  return %s new %s%s(%s);", pureVirtRetstr, pxyclsp, parent.Type().Spelling(), paramStr)
	}
	this.cp.APf("main", "}")
	this.genMethodFooter(clsctx, cursor, parent)
	if found && funco.Since != "" {
		this.cp.APf("main", "#endif // QT_VERSION >= %s", sinceVer2Hex(funco.Since))
	}
	this.cp.APf("main", "")
}

func (this *GenerateInlinev0) genDtor(clsctx *GenClassContext, cursor, parent clang.Cursor) {
	pparent := parent.SemanticParent()

	// this.cp.APf("main", "extern \"C\" Q_DECL_EXPORT")
	this.cp.APf("main", "void %s(void *this_) {", this.mangler.crc32p(cursor))
	if strings.HasPrefix(pparent.Spelling(), "Qt") {
		this.cp.APf("main", "  delete (%s::%s*)(this_);", pparent.Spelling(), parent.Spelling())
	} else {
		this.cp.APf("main", "  delete (%s*)(this_);", parent.Type().Spelling())
	}
	this.cp.APf("main", "}")
	this.genMethodFooter(clsctx, cursor, parent)
}

// 在该类没有显式的声明析构方法时，补充生成一个默认析构方法
func (this *GenerateInlinev0) genDtorNotsee(cursor, parent clang.Cursor) {
	// pparent := parent.SemanticParent()

	this.cp.APf("main", "")
	// this.cp.APf("main", "extern \"C\" Q_DECL_EXPORT")
	this.cp.APf("main", "void C_ZN%d%sD2Ev(void *this_) {", len(cursor.Spelling()), cursor.Spelling())
	if strings.HasPrefix(parent.Spelling(), "Qt") {
		this.cp.APf("main", "  delete (%s::%s*)(this_);", parent.Spelling(), cursor.Spelling())
	} else {
		this.cp.APf("main", "  delete (%s*)(this_);", cursor.Type().Spelling())
	}
	this.cp.APf("main", "}")
}

// for CXXMethodDecl, c++filt result endwiths &&
func (this *GenerateInlinev0) isrefrvalue(cursor clang.Cursor) bool {
	rety := cursor.Type()
	return rety.RefQualifier() == clang.RefQualifier_RValue
}

func (this *GenerateInlinev0) genNonStaticMethod(clsctx *GenClassContext, cursor, parent clang.Cursor, withParentSelector string) {
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	pparent := parent.SemanticParent()
	pparentstr := ""
	if strings.HasPrefix(pparent.Spelling(), "Qt") {
		// pparentstr = fmt.Sprintf("%s::", pparent.Spelling())
	}

	funco, found := (*parser.Function)(nil), false
	if clsctx.clso != nil {
		funco, found = qdi.findMethod(clsctx.clso, cursor)
	}

	if found && funco.Since != "" {
		this.cp.APf("main", "#if QT_VERSION >= %s", sinceVer2Hex(funco.Since))
	}

	if featname := is_feated_method(cursor); featname != "" {
		this.cp.APf("main", "#if QT_CONFIG(%s)", featname)
	}

	isrefrval := this.isrefrvalue(cursor) // if endwiths &&

	vaargstr := ""
	vaprmstr := ""
	if cursor.IsVariadic() && cursor.NumArguments() > 0 {
		for i := 0; i < 10; i++ {
			vaargstr += fmt.Sprintf(",void*a%d ", i)
			vaprmstr += fmt.Sprintf(",a%d ", i)
		}
	}
	// this.cp.APf("main", "extern \"C\" Q_DECL_EXPORT")
	// this.cp.APf("main", `__attribute__((visibility ("hidden")))`)
	this.cp.APf("main", "static")
	this.cp.APf("main", "void %s(%s%s) {", this.mangler.crc32p(cursor), argStr, vaargstr)
	this.cp.APf("main", "  ((%s%s*)0)->%s%s(%s%s);", pparentstr, parent.Type().Spelling(),
		withParentSelector, cursor.Spelling(), paramStr, vaprmstr)
	if isrefrval { // fix refqualifier rvalue
		this.cp.APf("main", "  (%s%s{}).%s%s(%s%s);", pparentstr, parent.Type().Spelling(),
			withParentSelector, cursor.Spelling(), paramStr, vaprmstr)
	}
	this.cp.APf("main", "  //%v auto x = %s;", isrefrval, this.mthFnptrDesc(cursor, parent))
	this.cp.APf("main", "}")
	this.genMethodFooter(clsctx, cursor, parent)

	if featname := is_feated_method(cursor); featname != "" {
		this.cp.APf("main", "#endif // QT_CONFIG(%s)", featname)
	}

	if found && funco.Since != "" {
		this.cp.APf("main", "#endif // QT_VERSION >= %s", sinceVer2Hex(funco.Since))
	}
	this.cp.APf("main", "")
}

func (this *GenerateInlinev0) genStaticMethod(clsctx *GenClassContext, cursor, parent clang.Cursor) {
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	pparent := parent.SemanticParent()
	pparentstr := ""
	if strings.HasPrefix(pparent.Spelling(), "Qt") {
		// pparentstr = fmt.Sprintf("%s::", pparent.Spelling())
	}

	funco, found := (*parser.Function)(nil), false
	if clsctx.clso != nil {
		funco, found = qdi.findMethod(clsctx.clso, cursor)
	}

	if found && funco.Since != "" {
		this.cp.APf("main", "#if QT_VERSION >= %s", sinceVer2Hex(funco.Since))
	}

	if featname := is_feated_method(cursor); featname != "" {
		this.cp.APf("main", "#if QT_CONFIG(%s)", featname)
	}

	vaargstr := ""
	vaprmstr := ""
	if cursor.IsVariadic() && cursor.NumArguments() > 0 {
		for i := 0; i < 10; i++ {
			vaargstr += fmt.Sprintf(",void*a%d ", i)
			vaprmstr += fmt.Sprintf(",a%d ", i)
		}
	}

	// this.cp.APf("main", "extern \"C\" Q_DECL_EXPORT")
	// this.cp.APf("main", `__attribute__((visibility ("hidden")))`)
	this.cp.APf("main", "static")
	this.cp.APf("main", "void %s(%s%s) {", this.mangler.crc32p(cursor), argStr, vaargstr)
	this.cp.APf("main", "  %s%s::%s(%s%s);",
		pparentstr, parent.Type().Spelling(), cursor.Spelling(), paramStr, vaprmstr)
	this.cp.APf("main", "}")
	this.genMethodFooter(clsctx, cursor, parent)

	if featname := is_feated_method(cursor); featname != "" {
		this.cp.APf("main", "#endif // QT_CONFIG(%s)", featname)
	}

	if found && funco.Since != "" {
		this.cp.APf("main", "#endif // QT_VERSION >= %s", sinceVer2Hex(funco.Since))
	}
	this.cp.APf("main", "")
}

func (this *GenerateInlinev0) genProtectedCallback(cursor, parent clang.Cursor) {
	// this.genMethodHeader(cursor, parent)
}

func (this *GenerateInlinev0) genArgsPxy(cursor, parent clang.Cursor) {
	this.argDesc = make([]string, 0)
	this.argtyDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArgPxy(argc, cursor, idx)
	}
	// log.Println(strings.Join(this.argDesc, ", "), this.mangler.convTo(cursor))
}

func (this *GenerateInlinev0) genArgPxy(cursor, parent clang.Cursor, idx int) {
	// log.Println(cursor.DisplayName(), cursor.Type().Spelling(), cursor.Type().Kind() == clang.Type_LValueReference, this.mangler.convTo(parent))
	csty := cursor.Type()
	argName := this.genParamRefName(cursor, parent, idx)
	if csty.Kind() == clang.Type_LValueReference {
		// 转成指针
	}
	if strings.Contains(csty.CanonicalType().Spelling(), "QFlags<") {
		this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s",
			csty.CanonicalType().Spelling(), argName))
		this.argtyDesc = append(this.argtyDesc, csty.CanonicalType().Spelling())
	} else {
		log.Println(cursor.Type().Kind(), csty.Spelling(), parent.SemanticParent().Spelling(), parent.DisplayName())
		if TypeIsFuncPointer(csty) {
			this.argDesc = append(this.argDesc,
				strings.Replace(csty.Spelling(), "(*)", fmt.Sprintf("(*%s)", argName), 1))
			this.argtyDesc = append(this.argtyDesc, cursor.Type().Spelling())
		} else if TypeIsCharPtrPtr(cursor.Type()) {
			this.argDesc = append(this.argDesc, fmt.Sprintf("char** %s", argName))
			this.argtyDesc = append(this.argtyDesc, "char**")
		} else if csty.Kind() == clang.Type_ConstantArray {
			log.Println(cursor.Spelling(), parent.Spelling(), csty.ArraySize(), csty.ArrayElementType().Spelling())
			elemty := csty.ArrayElementType()
			this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s[%d]", elemty.Spelling(), argName, csty.ArraySize()))
			this.argtyDesc = append(this.argtyDesc, fmt.Sprintf("%s[%d]", elemty.Spelling(), csty.ArraySize()))
		} else if csty.Kind() == clang.Type_IncompleteArray {
			log.Println(cursor.Spelling(), parent.Spelling(), csty.ArraySize(), csty.ArrayElementType().Spelling())
			elemty := csty.ArrayElementType()
			this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s[]", elemty.Spelling(), argName))
			this.argtyDesc = append(this.argtyDesc, fmt.Sprintf("%s[]", elemty.Spelling()))
		} else if (csty.Kind() == clang.Type_IncompleteArray ||
			csty.Kind() == clang.Type_ConstantArray) &&
			csty.ElementType().Kind() == clang.Type_Pointer {
			this.argDesc = append(this.argDesc, fmt.Sprintf("void** %s", argName))
			this.argtyDesc = append(this.argtyDesc, "void**")
		} else if csty.Kind() == clang.Type_IncompleteArray ||
			csty.Kind() == clang.Type_ConstantArray {
			this.argDesc = append(this.argDesc, fmt.Sprintf("void* %s", argName))
			this.argtyDesc = append(this.argtyDesc, "void*")
			// idx := strings.Index(cursor.Type().Spelling(), " [")
			// this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s %s",
			//	cursor.Type().Spelling()[0:idx], cursor.Spelling(), cursor.Type().Spelling()[idx+1:]))
		} else {
			this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s", csty.Spelling(), argName))
			this.argtyDesc = append(this.argtyDesc, csty.Spelling())
		}
	}
}

func (this *GenerateInlinev0) genArgs(cursor, parent clang.Cursor) {
	this.argDesc = make([]string, 0)
	this.argtyDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		argty := argc.Type()
		log.Println(parent.Spelling(), cursor.Spelling(),
			idx, argc.HasDefaultArg(), argty.Kind(), argty)
		isfuncptrty := argty.IsFunctionPointerType()
		unconstty := argty.RemoveLocalConst()
		log.Println(parent.Spelling(), cursor.Spelling(),
			idx, argc.HasDefaultArg(), isfuncptrty,
			argty.Spelling(), unconstty.Spelling())
		this.genArg(argc, cursor, idx)
	}
	// log.Println(strings.Join(this.argDesc, ", "), this.mangler.convTo(cursor))
}

func (this *GenerateInlinev0) genArg(cursor, parent clang.Cursor, idx int) {
	// log.Println(cursor.DisplayName(), cursor.Type().Spelling(), cursor.Type().Kind() == clang.Type_LValueReference, this.mangler.convTo(parent))
	// csty := cursor.Type()
	canty := cursor.Type().CanonicalType()
	argName := this.genParamRefName(cursor, parent, idx)

	tystr := canty.Spelling()
	if strings.Contains(tystr, "(*)(") || strings.Contains(tystr, "::*)(") ||
		strings.Contains(tystr, "__va_list_tag") {
		tystr = cursor.Type().Spelling()
	}
	argDesc := fmt.Sprintf("%s %s", tystr, argName)
	this.argDesc = append(this.argDesc, argDesc)
	argtyDesc := fmt.Sprintf("%s", cursor.Type().Spelling())
	this.argtyDesc = append(this.argtyDesc, argtyDesc)

}

func (this *GenerateInlinev0) genParamsPxy(cursor, parent clang.Cursor) {
	this.paramDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genParamPxy(argc, cursor, idx)
	}
}

func (this *GenerateInlinev0) genParamPxy(cursor, parent clang.Cursor, idx int) {
	csty := cursor.Type()
	forceConvStr := ""
	log.Println(csty.Kind().String(), csty.Spelling(), parent.Spelling(), csty.PointeeType().Kind().String(), csty.ArrayElementType().Kind().String())
	if csty.Kind() == clang.Type_Record { //} &&
		// (parent.Kind() != clang.Cursor_Constructor && !this.hasVirtualProjected) {
		// forceConvStr = "*"
	} else if TypeIsCharPtrPtr(csty) {
		// forceConvStr = "(char**)"
	}
	argName := this.genParamRefName(cursor, parent, idx)
	this.paramDesc = append(this.paramDesc, fmt.Sprintf("%s%s", forceConvStr, argName))
}

func (this *GenerateInlinev0) genParamsPxyCall(cursor, parent clang.Cursor) {
	this.paramDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genParamPxyCall(argc, cursor, idx)
	}
}

func (this *GenerateInlinev0) genParamPxyCall(cursor, parent clang.Cursor, idx int) {
	csty := cursor.Type()
	forceConvStr := ""
	log.Println(csty.Kind().String(), csty.Spelling(), parent.Spelling(), csty.PointeeType().Kind().String(), csty.ArrayElementType().Kind().String())
	if csty.Kind() == clang.Type_Record { //} &&
		// (parent.Kind() != clang.Cursor_Constructor && !this.hasVirtualProjected) {
		forceConvStr = fmt.Sprintf("(%s*)&", csty.Declaration().Spelling())
	} else if csty.Kind() == clang.Type_LValueReference &&
		csty.PointeeType().Kind() == clang.Type_Record {
		forceConvStr = fmt.Sprintf("(%s*)&", csty.PointeeType().Declaration().Spelling())
	} else if TypeIsCharPtrPtr(csty) {
		// forceConvStr = "(char**)"
	} else {
		log.Println()
	}

	argName := this.genParamRefName(cursor, parent, idx)
	this.paramDesc = append(this.paramDesc, fmt.Sprintf("%s%s", forceConvStr, argName))
}

func (this *GenerateInlinev0) genParams(cursor, parent clang.Cursor) {
	this.paramDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genParam(argc, cursor, idx)
	}
}

func (this *GenerateInlinev0) genParam(cursor, parent clang.Cursor, idx int) {
	csty := cursor.Type()
	forceConvStr := ""
	log.Println(csty.Kind().String(), csty.Spelling(), parent.Spelling(), csty.PointeeType().Kind().String(), csty.ArrayElementType().Kind().String())

	argName := this.genParamRefName(cursor, parent, idx)
	if idx == 0 && // for argc
		funk.ContainsString([]string{"QCoreApplication", "QGuiApplication", "QApplication", "QAndroidService"}, parent.Spelling()) {
		this.paramDesc = append(this.paramDesc, fmt.Sprintf("*(new int(%s))", argName))
	} else {
		this.paramDesc = append(this.paramDesc, fmt.Sprintf("%s%s", forceConvStr, argName))
	}
}

func (this *GenerateInlinev0) genRet(cursor, parent clang.Cursor, idx int) {

}

//
func (this *GenerateInlinev0) genCSignature(cursor, parent clang.Cursor, idx int) {

}

func (this *GenerateInlinev0) genEnumsGlobal(cursor, parent clang.Cursor) {

}
func (this *GenerateInlinev0) genEnum() {

}

func (this *GenerateInlinev0) genFunctions(cursor, parent clang.Cursor) {
	// this.genHeader(cursor, parent)
	skipKeys := []string{"QKeySequence", "QVector2D", "QPointingDeviceUniqueId", "QFont", "QMatrix",
		"QTransform", "QPixelFormat", "QRawFont", "QVector3D", "QVector4D",
		"QOpenGLVersionStatus", "QOpenGLVersionProfile"}
	hasSkipKey := func(c clang.Cursor) bool {
		for _, k := range skipKeys {
			if strings.Contains(c.DisplayName(), k) {
				return true
			}
		}
		return false
	}

	grfuncs := this.groupFunctionsByModule()
	qtmods := []string{}
	for qtmod, _ := range grfuncs {
		qtmods = append(qtmods, qtmod)
	}
	sort.Strings(qtmods)

	for _, qtmod := range qtmods {
		funcs := grfuncs[qtmod]
		log.Println(qtmod, len(funcs))
		this.cp = NewCodePager()
		// write code
		for _, mod := range modDeps[qtmod] {
			incmod := getIncludeNameByModule(mod)
			this.cp.APf("header", "#include <Qt%s>", incmod)
		}
		this.cp.APf("header", "#include <Qt%s>", getIncludeNameByModule(qtmod))
		this.cp.APf("header", "#include \"hidden_symbols.h\"")

		sort.Slice(funcs, func(i int, j int) bool {
			return funcs[i].Mangling() > funcs[j].Mangling()
		})
		for _, fc := range funcs {
			log.Println(fc.Kind(), fc.Spelling(), fc.Mangling(), fc.DisplayName(), fc.IsCursorDefinition(), fc.SemanticParent().Spelling(), fc.Type().Spelling(), fc.SemanticParent().Kind().String(), fc.NumTemplateArguments())
			if !is_qt_global_func(fc) {
				log.Println("skip global function ", fc.Spelling(), fc.IsCursorDefinition(), this.mangler.origin(fc))
				continue
			}

			if fc.NumTemplateArguments() > 0 {
				log.Println("skip global function ", fc.Spelling(),
					fc.NumTemplateArguments(), fc.TemplateArgumentType(0).Spelling(),
					fc.TemplateArgumentKind(0).Spelling(), fc.TemplateCursorKind().Spelling())
				continue
			}
			if strings.ContainsAny(fc.DisplayName(), "<>") {
				// log.Println("skip global function ", fc.Spelling())
				// continue
			}
			if strings.Contains(fc.DisplayName(), "Rgba64") {
				// log.Println("skip global function ", fc.Spelling())
				// continue
			}
			if strings.Contains(fc.ResultType().Spelling(), "Rgba64") {
				// log.Println("skip global function ", fc.Spelling())
				// continue
			}
			if hasSkipKey(fc) {
				//log.Println("skip global function ", fc.Spelling())
				//continue
			}

			if this.filter.skipFunc(fc) {
				log.Println("skip global function ", fc.Spelling())
				continue
			}

			if _, ok := this.funcMangles[fc.Spelling()]; ok {
				this.funcMangles[fc.Spelling()] += 1
			} else {
				this.funcMangles[fc.Spelling()] = 0
			}
			olidx := this.funcMangles[fc.Spelling()]
			this.genFunction(fc, olidx)
		}

		this.saveCodeToFile(qtmod, "qfunctions")
	}
}

func (this *GenerateInlinev0) genFunction(cursor clang.Cursor, olidx int) {
	log.Println(cursor.DisplayName(), len(this.funcs))
	this.genArgs(cursor, cursor.SemanticParent())
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, cursor.SemanticParent())
	paramStr := strings.Join(this.paramDesc, ", ")

	if strings.Contains(argStr, "DropActions::enum_type") {
		// log.Fatalln(cursor.DisplayName(), cursor.Spelling(), argStr, paramStr)
	}

	parent := cursor.SemanticParent()
	nsfix := gopp.IfElseStr(parent.Kind() == clang.Cursor_Namespace, parent.Spelling()+"::", "")

	hasLongDoubleArg := FuncHasLongDoubleArg(cursor)
	if hasLongDoubleArg {
		this.cp.APf("main", "#ifndef Q_OS_DARWIN")
	}
	overloadSuffix := gopp.IfElseStr(olidx == 0, "", fmt.Sprintf("_%d", olidx))
	clsctx := &GenClassContext{}
	this.genMethodHeader(clsctx, cursor, cursor.SemanticParent())

	vaargstr := ""
	vaprmstr := ""
	if cursor.IsVariadic() && cursor.NumArguments() > 0 {
		for i := 0; i < 10; i++ {
			vaargstr += fmt.Sprintf(",void*a%d ", i)
			vaprmstr += fmt.Sprintf(",a%d ", i)
		}
	}

	// this.cp.APf("main", "extern \"C\" Q_DECL_EXPORT")
	this.cp.APf("main", "void %s%s(%s%s) {",
		this.mangler.crc32p(cursor), overloadSuffix, argStr, vaargstr)
	this.cp.APf("main", "  %s%s(%s%s);", nsfix, cursor.Spelling(), paramStr, vaprmstr)
	this.cp.APf("main", "}")
	this.genMethodFooter(nil, cursor, cursor)
	if hasLongDoubleArg {
		this.cp.APf("main", "#endif")
	}
	this.cp.APf("main", "")
}

func (this *GenerateInlinev0) genConstantsGlobal(cursor, parent clang.Cursor) {
}

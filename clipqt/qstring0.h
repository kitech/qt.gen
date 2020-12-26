#include <QtCore>

class xQByteArray : public QString {
 public:
    xQByteArray(const char *, int size = -1);
    char *data();
    inline int size() const;
};

class xQString : public QString {
 public:
    xQString(const char *ch);

    int length() const;
    QByteArray toUtf8() const;
    static inline QString fromUtf8(const char *str, int size = -1);
};

class xQVariant : public QVariant {
public:
    xQVariant() noexcept ;
//     ~QVariant();
//     QVariant(Type type);
//     QVariant(int typeId, const void *copy);
//     QVariant(int typeId, const void *copy, uint flags);
//     QVariant(const QVariant &other);

// #ifndef QT_NO_DATASTREAM
//     QVariant(QDataStream &s);
// #endif

//     QVariant(int i);
//     QVariant(uint ui);
//     QVariant(qlonglong ll);
//     QVariant(qulonglong ull);
//     QVariant(bool b);
//     QVariant(double d);
//     QVariant(float f);
// #ifndef QT_NO_CAST_FROM_ASCII
//     QT_ASCII_CAST_WARN QVariant(const char *str);
// #endif

//     QVariant(const QByteArray &bytearray);
//     QVariant(const QBitArray &bitarray);
//     QVariant(const QString &string);
//     QVariant(QLatin1String string);
//     QVariant(const QStringList &stringlist);
//     QVariant(QChar qchar);
//     QVariant(const QDate &date);
//     QVariant(const QTime &time);
//     QVariant(const QDateTime &datetime);
//     QVariant(const QList<QVariant> &list);
//     QVariant(const QMap<QString,QVariant> &map);
//     QVariant(const QHash<QString,QVariant> &hash);
// #ifndef QT_NO_GEOM_VARIANT
//     QVariant(const QSize &size);
//     QVariant(const QSizeF &size);
//     QVariant(const QPoint &pt);
//     QVariant(const QPointF &pt);
//     QVariant(const QLine &line);
//     QVariant(const QLineF &line);
//     QVariant(const QRect &rect);
//     QVariant(const QRectF &rect);
// #endif
//     QVariant(const QLocale &locale);
// #ifndef QT_NO_REGEXP
//     QVariant(const QRegExp &regExp);
// #endif // QT_NO_REGEXP
// #if QT_CONFIG(regularexpression)
//     QVariant(const QRegularExpression &re);
// #endif // QT_CONFIG(regularexpression)
// #if QT_CONFIG(easingcurve)
//     QVariant(const QEasingCurve &easing);
// #endif
//     QVariant(const QUuid &uuid);
// #ifndef QT_BOOTSTRAPPED
//     QVariant(const QUrl &url);
//     QVariant(const QJsonValue &jsonValue);
//     QVariant(const QJsonObject &jsonObject);
//     QVariant(const QJsonArray &jsonArray);
//     QVariant(const QJsonDocument &jsonDocument);
// #endif // QT_BOOTSTRAPPED
// #if QT_CONFIG(itemmodel)
//     QVariant(const QModelIndex &modelIndex);
//     QVariant(const QPersistentModelIndex &modelIndex);
// #endif

//     QVariant& operator=(const QVariant &other);
//     inline QVariant(QVariant &&other) noexcept : d(other.d)
//     { other.d = Private(); }
//     inline QVariant &operator=(QVariant &&other) noexcept
//     { qSwap(d, other.d); return *this; }

//     inline void swap(QVariant &other) noexcept { qSwap(d, other.d); }

//     Type type() const;
//     int userType() const;
//     const char *typeName() const;

//     bool canConvert(int targetTypeId) const;
//     bool convert(int targetTypeId);

//     inline bool isValid() const;
//     bool isNull() const;

//     void clear();

//     void detach();
//     inline bool isDetached() const;

//     int toInt(bool *ok = nullptr) const;
//     uint toUInt(bool *ok = nullptr) const;
//     qlonglong toLongLong(bool *ok = nullptr) const;
//     qulonglong toULongLong(bool *ok = nullptr) const;
//     bool toBool() const;
//     double toDouble(bool *ok = nullptr) const;
//     float toFloat(bool *ok = nullptr) const;
//     qreal toReal(bool *ok = nullptr) const;
//     QByteArray toByteArray() const;
//     QBitArray toBitArray() const;
//     QString toString() const;
//     QStringList toStringList() const;
//     QChar toChar() const;
//     QDate toDate() const;
//     QTime toTime() const;
//     QDateTime toDateTime() const;
//     QList<QVariant> toList() const;
//     QMap<QString, QVariant> toMap() const;
//     QHash<QString, QVariant> toHash() const;

// #ifndef QT_NO_GEOM_VARIANT
//     QPoint toPoint() const;
//     QPointF toPointF() const;
//     QRect toRect() const;
//     QSize toSize() const;
//     QSizeF toSizeF() const;
//     QLine toLine() const;
//     QLineF toLineF() const;
//     QRectF toRectF() const;
// #endif
//     QLocale toLocale() const;
// #ifndef QT_NO_REGEXP
//     QRegExp toRegExp() const;
// #endif // QT_NO_REGEXP
// #if QT_CONFIG(regularexpression)
//     QRegularExpression toRegularExpression() const;
// #endif // QT_CONFIG(regularexpression)
// #if QT_CONFIG(easingcurve)
//     QEasingCurve toEasingCurve() const;
// #endif
//     QUuid toUuid() const;
// #ifndef QT_BOOTSTRAPPED
//     QUrl toUrl() const;
//     QJsonValue toJsonValue() const;
//     QJsonObject toJsonObject() const;
//     QJsonArray toJsonArray() const;
//     QJsonDocument toJsonDocument() const;
// #endif // QT_BOOTSTRAPPED
// #if QT_CONFIG(itemmodel)
//     QModelIndex toModelIndex() const;
//     QPersistentModelIndex toPersistentModelIndex() const;
// #endif

// #ifndef QT_NO_DATASTREAM
//     void load(QDataStream &ds);
//     void save(QDataStream &ds) const;
// #endif
//     static const char *typeToName(int typeId);
//     static Type nameToType(const char *name);

//     void *data();
//     const void *constData() const;
//     inline const void *data() const { return constData(); }

//     template<typename T>
//     inline void setValue(const T &value);

//     template<typename T>
//     inline T value() const
//     { return qvariant_cast<T>(*this); }

//     template<typename T>
//     static inline QVariant fromValue(const T &value)
//     { return QVariant(qMetaTypeId<T>(), &value, QTypeInfo<T>::isPointer); }

// #if (__has_include(<variant>) && __cplusplus >= 201703L) || defined(Q_CLANG_QDOC)
//     template<typename... Types>
//     static inline QVariant fromStdVariant(const std::variant<Types...> &value)
//     {
//         if (value.valueless_by_exception())
//             return QVariant();
//         return std::visit([](const auto &arg) { return fromValue(arg); }, value);
//     }
// #endif

//     template<typename T>
//     bool canConvert() const
//     { return canConvert(qMetaTypeId<T>()); }

//  public:
//     struct PrivateShared
//     {
//         inline PrivateShared(void *v) : ptr(v), ref(1) { }
//         void *ptr;
//         QAtomicInt ref;
//     };
//     struct Private
//     {
//         inline Private() noexcept : type(Invalid), is_shared(false), is_null(true)
//         { data.ptr = nullptr; }

//         // Internal constructor for initialized variants.
//         explicit inline Private(uint variantType) noexcept
//             : type(variantType), is_shared(false), is_null(false)
//         {}

// #if QT_VERSION < QT_VERSION_CHECK(6, 0, 0)
//         Private(const Private &other) noexcept
//             : data(other.data), type(other.type),
//               is_shared(other.is_shared), is_null(other.is_null)
//         {}
//         Private &operator=(const Private &other) noexcept = default;
// #endif
//         union Data
//         {
//             char c;
//             uchar uc;
//             short s;
//             signed char sc;
//             ushort us;
//             int i;
//             uint u;
//             long l;
//             ulong ul;
//             bool b;
//             double d;
//             float f;
//             qreal real;
//             qlonglong ll;
//             qulonglong ull;
//             QObject *o;
//             void *ptr;
//             PrivateShared *shared;
//         } data;
//         uint type : 30;
//         uint is_shared : 1;
//         uint is_null : 1;
//     };
//  public:
//     typedef void (*f_construct)(Private *, const void *);
//     typedef void (*f_clear)(Private *);
//     typedef bool (*f_null)(const Private *);
// #ifndef QT_NO_DATASTREAM
//     typedef void (*f_load)(Private *, QDataStream &);
//     typedef void (*f_save)(const Private *, QDataStream &);
// #endif
//     typedef bool (*f_compare)(const Private *, const Private *);
//     typedef bool (*f_convert)(const QVariant::Private *d, int t, void *, bool *);
//     typedef bool (*f_canConvert)(const QVariant::Private *d, int t);
//     typedef void (*f_debugStream)(QDebug, const QVariant &);
//     struct Handler {
//         f_construct construct;
//         f_clear clear;
//         f_null isNull;
// #ifndef QT_NO_DATASTREAM
//         f_load load;
//         f_save save;
// #endif
//         f_compare compare;
//         f_convert convert;
//         f_canConvert canConvert;
//         f_debugStream debugStream;
//     };

//     inline bool operator==(const QVariant &v) const
//     { return cmp(v); }
//     inline bool operator!=(const QVariant &v) const
//     { return !cmp(v); }
// #if QT_DEPRECATED_SINCE(5, 15)
//     QT_DEPRECATED inline bool operator<(const QVariant &v) const
//     { return compare(v) < 0; }
//     QT_DEPRECATED inline bool operator<=(const QVariant &v) const
//     { return compare(v) <= 0; }
//     QT_DEPRECATED inline bool operator>(const QVariant &v) const
//     { return compare(v) > 0; }
//     QT_DEPRECATED inline bool operator>=(const QVariant &v) const
//     { return compare(v) >= 0; }
// #endif

}; // QVariant

class xQUrl : public QUrl {
public:
    xQUrl();
    explicit xQUrl(const QString &url, ParsingMode mode = TolerantMode);

    void setUrl(const QString &url, ParsingMode mode = TolerantMode);
    QString url(FormattingOptions options = FormattingOptions(PrettyDecoded)) const;
    QString toString(FormattingOptions options = FormattingOptions(PrettyDecoded)) const;
    QString toDisplayString(FormattingOptions options = FormattingOptions(PrettyDecoded)) const;
}; // QUrl

class xQSize : public QSize
{
 public:
    xQSize(int w, int h);

    inline int width() const noexcept;
    inline int height() const noexcept;

    inline int &rwidth() noexcept;
    inline int &rheight() noexcept;
};


class xQEvent : public QEvent            // event base class
{
public:
    // explicit QEvent(Type type);
    // QEvent(const QEvent &other);
    // virtual ~QEvent();
    // QEvent &operator=(const QEvent &other);
    inline Type type() const { return static_cast<Type>(t); }
    inline bool spontaneous() const { return spont; }

    inline void setAccepted(bool accepted) { m_accept = accepted; }
    inline bool isAccepted() const { return m_accept; }

    inline void accept() { m_accept = true; }
    inline void ignore() { m_accept = false; }

    static int registerEventType(int hint = -1) noexcept;
};

class xQTimerEvent : public QTimerEvent
{
public:
    // explicit QTimerEvent( int timerId );
    // ~QTimerEvent();
    // int timerId() const { return id; }
};

class xQChildEvent : public QChildEvent
{
public:
    // QChildEvent( Type type, QObject *child );
    // ~QChildEvent();
    // QObject *child() const { return c; }
    // bool added() const { return type() == ChildAdded; }
    // bool polished() const { return type() == ChildPolished; }
    // bool removed() const { return type() == ChildRemoved; }
};

class xQDynamicPropertyChangeEvent : public QDynamicPropertyChangeEvent
{
public:
    // explicit QDynamicPropertyChangeEvent(const QByteArray &name);
    // ~QDynamicPropertyChangeEvent();

    // inline QByteArray propertyName() const { return n; }
};

class xQDeferredDeleteEvent : public QDeferredDeleteEvent
{
public:
    // explicit QDeferredDeleteEvent();
    // ~QDeferredDeleteEvent();
    // int loopLevel() const { return level; }
};


class xQMetaClassInfo : public QMetaClassInfo {
public:
};

class xQMetaEnum : public QMetaEnum {
public:
};

class xQMetaProperty : public QMetaProperty {
public:
};

class xQMetaMethod : public QMetaMethod
{
 public:
    QByteArray methodSignature() const;
    QByteArray name() const;
    const char *typeName() const;
    int returnType() const;
    int parameterCount() const;
    int parameterType(int index) const;
    void getParameterTypes(int *types) const;
    // QList<QByteArray> parameterTypes() const;
    // QList<QByteArray> parameterNames() const;
    const char *tag() const;
    Access access() const;
    MethodType methodType() const;
    int attributes() const;
    int methodIndex() const;
    int revision() const;

    inline const QMetaObject *enclosingMetaObject() const;
};

class xQMetaObject: public QMetaObject
{
 public:
    const char *className() const;
    const QMetaObject *superClass() const;

    bool inherits(const QMetaObject *metaObject) const noexcept;
    QObject *cast(QObject *obj) const;
    const QObject *cast(const QObject *obj) const;

    int methodOffset() const;
    int enumeratorOffset() const;
    int propertyOffset() const;
    int classInfoOffset() const;

    int constructorCount() const;
    int methodCount() const;
    int enumeratorCount() const;
    int propertyCount() const;
    int classInfoCount() const;

    int indexOfConstructor(const char *constructor) const;
    int indexOfMethod(const char *method) const;
    int indexOfSignal(const char *signal) const;
    int indexOfSlot(const char *slot) const;
    int indexOfEnumerator(const char *name) const;
    int indexOfProperty(const char *name) const;
    int indexOfClassInfo(const char *name) const;

    QMetaMethod constructor(int index) const;
    QMetaMethod method(int index) const;
    QMetaEnum enumerator(int index) const;
    QMetaProperty property(int index) const;
    QMetaClassInfo classInfo(int index) const;
    QMetaProperty userProperty() const;

};

class xQGenericReturnArgument : public QGenericReturnArgument {
public:
};

class xQGenericArgument : public QGenericArgument {
public:
}

class xQObject : public QObject
{
 public:

    QString objectName() const;
    void setObjectName(const QString &name);

    inline bool isWidgetType() const;
    inline bool isWindowType() const ;

    inline bool signalsBlocked() const noexcept ;
    bool blockSignals(bool b) noexcept;

    QThread *thread() const;
    void moveToThread(QThread *thread);

    void setParent(QObject *parent);
};


class xQThread : public QThread {
public:
};

class xQRect : public QRect {
public:
};

class xQRectF : public QRectF {
public:
};

class xQPoint : public QPoint {
public:
};

class xQModelIndex : public QModelIndex
{
public:
    Q_DECL_CONSTEXPR inline xQModelIndex() noexcept : r(-1), c(-1), i(0), m(nullptr) {}
    // compiler-generated copy/move ctors/assignment operators are fine!
    Q_DECL_CONSTEXPR inline int row() const noexcept;
    Q_DECL_CONSTEXPR inline int column() const noexcept ;
    Q_DECL_CONSTEXPR inline quintptr internalId() const noexcept;
    inline void *internalPointer() const noexcept;
    inline QModelIndex parent() const;
    inline QModelIndex sibling(int row, int column) const;
    inline QModelIndex siblingAtColumn(int column) const;
    inline QModelIndex siblingAtRow(int row) const;
// #if QT_DEPRECATED_SINCE(5, 8)
//     QT_DEPRECATED_X("Use QAbstractItemModel::index") inline QModelIndex child(int row, int column) const;
// #endif
//     inline QVariant data(int role = Qt::DisplayRole) const;
//     inline Qt::ItemFlags flags() const;
//     Q_DECL_CONSTEXPR inline const QAbstractItemModel *model() const noexcept { return m; }
//     Q_DECL_CONSTEXPR inline bool isValid() const noexcept { return (r >= 0) && (c >= 0) && (m != nullptr); }
//     Q_DECL_CONSTEXPR inline bool operator==(const QModelIndex &other) const noexcept
//         { return (other.r == r) && (other.i == i) && (other.c == c) && (other.m == m); }
//     Q_DECL_CONSTEXPR inline bool operator!=(const QModelIndex &other) const noexcept
//         { return !(*this == other); }
//     Q_DECL_CONSTEXPR inline bool operator<(const QModelIndex &other) const noexcept
//         {
//             return  r <  other.r
//                 || (r == other.r && (c <  other.c
//                                  || (c == other.c && (i <  other.i
//                                                   || (i == other.i && std::less<const QAbstractItemModel *>()(m, other.m))))));
//         }
}; // QModelIndex

class  xQPersistentModelIndex: public QPersistentModelIndex
{
public:
    xQPersistentModelIndex();
//     QPersistentModelIndex(const QModelIndex &index);
//     QPersistentModelIndex(const QPersistentModelIndex &other);
//     ~QPersistentModelIndex();
//     bool operator<(const QPersistentModelIndex &other) const;
//     bool operator==(const QPersistentModelIndex &other) const;
//     inline bool operator!=(const QPersistentModelIndex &other) const
//     { return !operator==(other); }
//     QPersistentModelIndex &operator=(const QPersistentModelIndex &other);
//     inline QPersistentModelIndex(QPersistentModelIndex &&other) noexcept
//         : d(other.d) { other.d = nullptr; }
//     inline QPersistentModelIndex &operator=(QPersistentModelIndex &&other) noexcept
//     { qSwap(d, other.d); return *this; }
//     inline void swap(QPersistentModelIndex &other) noexcept { qSwap(d, other.d); }
//     bool operator==(const QModelIndex &other) const;
//     bool operator!=(const QModelIndex &other) const;
//     QPersistentModelIndex &operator=(const QModelIndex &other);
//     operator const QModelIndex&() const;
//     int row() const;
//     int column() const;
//     void *internalPointer() const;
//     quintptr internalId() const;
//     QModelIndex parent() const;
//     QModelIndex sibling(int row, int column) const;
// #if QT_DEPRECATED_SINCE(5, 8)
//     QT_DEPRECATED_X("Use QAbstractItemModel::index") QModelIndex child(int row, int column) const;
// #endif
//     QVariant data(int role = Qt::DisplayRole) const;
//     Qt::ItemFlags flags() const;
//     const QAbstractItemModel *model() const;
//     bool isValid() const;
}; // QModelIndex

class xQAbstractItemModel : public QAbstractItemModel
{
public:

    explicit xQAbstractItemModel(QObject *parent = nullptr);
//     virtual ~QAbstractItemModel();

//     Q_INVOKABLE bool hasIndex(int row, int column, const QModelIndex &parent = QModelIndex()) const;
//     Q_INVOKABLE virtual QModelIndex index(int row, int column,
//                               const QModelIndex &parent = QModelIndex()) const = 0;
//     Q_INVOKABLE virtual QModelIndex parent(const QModelIndex &child) const = 0;

//     Q_INVOKABLE virtual QModelIndex sibling(int row, int column, const QModelIndex &idx) const;
//     Q_INVOKABLE virtual int rowCount(const QModelIndex &parent = QModelIndex()) const = 0;
//     Q_INVOKABLE virtual int columnCount(const QModelIndex &parent = QModelIndex()) const = 0;
//     Q_INVOKABLE virtual bool hasChildren(const QModelIndex &parent = QModelIndex()) const;

//     Q_INVOKABLE virtual QVariant data(const QModelIndex &index, int role = Qt::DisplayRole) const = 0;
//     Q_INVOKABLE virtual bool setData(const QModelIndex &index, const QVariant &value, int role = Qt::EditRole);

//     Q_INVOKABLE virtual QVariant headerData(int section, Qt::Orientation orientation,
//                                 int role = Qt::DisplayRole) const;
//     virtual bool setHeaderData(int section, Qt::Orientation orientation, const QVariant &value,
//                                int role = Qt::EditRole);

//     virtual QMap<int, QVariant> itemData(const QModelIndex &index) const;
//     virtual bool setItemData(const QModelIndex &index, const QMap<int, QVariant> &roles);
// #if QT_VERSION >= QT_VERSION_CHECK(6, 0, 0)
//     virtual bool clearItemData(const QModelIndex &index);
// #endif

//     virtual QStringList mimeTypes() const;
//     virtual QMimeData *mimeData(const QModelIndexList &indexes) const;
//     virtual bool canDropMimeData(const QMimeData *data, Qt::DropAction action,
//                                  int row, int column, const QModelIndex &parent) const;
//     virtual bool dropMimeData(const QMimeData *data, Qt::DropAction action,
//                               int row, int column, const QModelIndex &parent);
//     virtual Qt::DropActions supportedDropActions() const;

//     virtual Qt::DropActions supportedDragActions() const;
// #if QT_DEPRECATED_SINCE(5, 0)
//     QT_DEPRECATED void setSupportedDragActions(Qt::DropActions actions)
//     { doSetSupportedDragActions(actions); }
// #endif

//     virtual bool insertRows(int row, int count, const QModelIndex &parent = QModelIndex());
//     virtual bool insertColumns(int column, int count, const QModelIndex &parent = QModelIndex());
//     virtual bool removeRows(int row, int count, const QModelIndex &parent = QModelIndex());
//     virtual bool removeColumns(int column, int count, const QModelIndex &parent = QModelIndex());
//     virtual bool moveRows(const QModelIndex &sourceParent, int sourceRow, int count,
//                           const QModelIndex &destinationParent, int destinationChild);
//     virtual bool moveColumns(const QModelIndex &sourceParent, int sourceColumn, int count,
//                              const QModelIndex &destinationParent, int destinationChild);

//     inline bool insertRow(int row, const QModelIndex &parent = QModelIndex());
//     inline bool insertColumn(int column, const QModelIndex &parent = QModelIndex());
//     inline bool removeRow(int row, const QModelIndex &parent = QModelIndex());
//     inline bool removeColumn(int column, const QModelIndex &parent = QModelIndex());
//     inline bool moveRow(const QModelIndex &sourceParent, int sourceRow,
//                         const QModelIndex &destinationParent, int destinationChild);
//     inline bool moveColumn(const QModelIndex &sourceParent, int sourceColumn,
//                            const QModelIndex &destinationParent, int destinationChild);

//     Q_INVOKABLE virtual void fetchMore(const QModelIndex &parent);
//     Q_INVOKABLE virtual bool canFetchMore(const QModelIndex &parent) const;
//     Q_INVOKABLE virtual Qt::ItemFlags flags(const QModelIndex &index) const;
//     virtual void sort(int column, Qt::SortOrder order = Qt::AscendingOrder);
//     virtual QModelIndex buddy(const QModelIndex &index) const;
//     Q_INVOKABLE virtual QModelIndexList match(const QModelIndex &start, int role,
//                                               const QVariant &value, int hits = 1,
//                                               Qt::MatchFlags flags =
//                                               Qt::MatchFlags(Qt::MatchStartsWith|Qt::MatchWrap)) const;
//     virtual QSize span(const QModelIndex &index) const;

//     virtual QHash<int,QByteArray> roleNames() const;

//     using QObject::parent;

//     enum LayoutChangeHint
//     {
//         NoLayoutChangeHint,
//         VerticalSortHint,
//         HorizontalSortHint
//     };
//     Q_ENUM(LayoutChangeHint)

//     enum class CheckIndexOption {
//         NoOption         = 0x0000,
//         IndexIsValid     = 0x0001,
//         DoNotUseParent   = 0x0002,
//         ParentIsInvalid  = 0x0004,
//     };
//     Q_ENUM(CheckIndexOption)
//     Q_DECLARE_FLAGS(CheckIndexOptions, CheckIndexOption)

//     Q_REQUIRED_RESULT bool checkIndex(const QModelIndex &index, CheckIndexOptions options = CheckIndexOption::NoOption) const;

// Q_SIGNALS:
//     void dataChanged(const QModelIndex &topLeft, const QModelIndex &bottomRight, const QVector<int> &roles = QVector<int>());
//     void headerDataChanged(Qt::Orientation orientation, int first, int last);
//     void layoutChanged(const QList<QPersistentModelIndex> &parents = QList<QPersistentModelIndex>(), QAbstractItemModel::LayoutChangeHint hint = QAbstractItemModel::NoLayoutChangeHint);
//     void layoutAboutToBeChanged(const QList<QPersistentModelIndex> &parents = QList<QPersistentModelIndex>(), QAbstractItemModel::LayoutChangeHint hint = QAbstractItemModel::NoLayoutChangeHint);

//     void rowsAboutToBeInserted(const QModelIndex &parent, int first, int last, QPrivateSignal);
//     void rowsInserted(const QModelIndex &parent, int first, int last, QPrivateSignal);

//     void rowsAboutToBeRemoved(const QModelIndex &parent, int first, int last, QPrivateSignal);
//     void rowsRemoved(const QModelIndex &parent, int first, int last, QPrivateSignal);

//     void columnsAboutToBeInserted(const QModelIndex &parent, int first, int last, QPrivateSignal);
//     void columnsInserted(const QModelIndex &parent, int first, int last, QPrivateSignal);

//     void columnsAboutToBeRemoved(const QModelIndex &parent, int first, int last, QPrivateSignal);
//     void columnsRemoved(const QModelIndex &parent, int first, int last, QPrivateSignal);

//     void modelAboutToBeReset(QPrivateSignal);
//     void modelReset(QPrivateSignal);

//     void rowsAboutToBeMoved( const QModelIndex &sourceParent, int sourceStart, int sourceEnd, const QModelIndex &destinationParent, int destinationRow, QPrivateSignal);
//     void rowsMoved( const QModelIndex &parent, int start, int end, const QModelIndex &destination, int row, QPrivateSignal);

//     void columnsAboutToBeMoved( const QModelIndex &sourceParent, int sourceStart, int sourceEnd, const QModelIndex &destinationParent, int destinationColumn, QPrivateSignal);
//     void columnsMoved( const QModelIndex &parent, int start, int end, const QModelIndex &destination, int column, QPrivateSignal);

// public Q_SLOTS:
//     virtual bool submit();
//     virtual void revert();

// protected Q_SLOTS:
// #if QT_VERSION >= QT_VERSION_CHECK(6, 0, 0)
//     virtual
// #endif
//     void resetInternalData();

// protected:
//     QAbstractItemModel(QAbstractItemModelPrivate &dd, QObject *parent = nullptr);

//     inline QModelIndex createIndex(int row, int column, void *data = nullptr) const;
//     inline QModelIndex createIndex(int row, int column, quintptr id) const;

//     void encodeData(const QModelIndexList &indexes, QDataStream &stream) const;
//     bool decodeData(int row, int column, const QModelIndex &parent, QDataStream &stream);

//     void beginInsertRows(const QModelIndex &parent, int first, int last);
//     void endInsertRows();

//     void beginRemoveRows(const QModelIndex &parent, int first, int last);
//     void endRemoveRows();

//     bool beginMoveRows(const QModelIndex &sourceParent, int sourceFirst, int sourceLast, const QModelIndex &destinationParent, int destinationRow);
//     void endMoveRows();

//     void beginInsertColumns(const QModelIndex &parent, int first, int last);
//     void endInsertColumns();

//     void beginRemoveColumns(const QModelIndex &parent, int first, int last);
//     void endRemoveColumns();

//     bool beginMoveColumns(const QModelIndex &sourceParent, int sourceFirst, int sourceLast, const QModelIndex &destinationParent, int destinationColumn);
//     void endMoveColumns();


// #if QT_DEPRECATED_SINCE(5,0)
//     QT_DEPRECATED void reset()
//     {
//         beginResetModel();
//         endResetModel();
//     }
// #endif

//     void beginResetModel();
//     void endResetModel();

//     void changePersistentIndex(const QModelIndex &from, const QModelIndex &to);
//     void changePersistentIndexList(const QModelIndexList &from, const QModelIndexList &to);
//     QModelIndexList persistentIndexList() const;

// #if QT_DEPRECATED_SINCE(5,0)
//     QT_DEPRECATED void setRoleNames(const QHash<int,QByteArray> &theRoleNames)
//     {
//         doSetRoleNames(theRoleNames);
//     }
// #endif

}; // QAbstractItemModel

class  xQAbstractTableModel : public QAbstractTableModel
{
 public:
    explicit xQAbstractTableModel(QObject *parent = nullptr);
    // ~QAbstractTableModel();

    // QModelIndex index(int row, int column, const QModelIndex &parent = QModelIndex()) const override;
    // QModelIndex sibling(int row, int column, const QModelIndex &idx) const override;
    // bool dropMimeData(const QMimeData *data, Qt::DropAction action,
    //                   int row, int column, const QModelIndex &parent) override;

    // Qt::ItemFlags flags(const QModelIndex &index) const override;

    // using QObject::parent;

}; // QAbstractTableModel

class xQAbstractListModel : public QAbstractListModel
{
 public:
    explicit xQAbstractListModel(QObject *parent = nullptr);
    // ~QAbstractListModel();

    // QModelIndex index(int row, int column = 0, const QModelIndex &parent = QModelIndex()) const override;
    // QModelIndex sibling(int row, int column, const QModelIndex &idx) const override;
    // bool dropMimeData(const QMimeData *data, Qt::DropAction action,
    //                   int row, int column, const QModelIndex &parent) override;

    // Qt::ItemFlags flags(const QModelIndex &index) const override;

    // using QObject::parent;

}; // QAbstractListModel


class xQCoreApplication : public QCoreApplication {
public:
//     enum { ApplicationFlags = QT_VERSION
//     };

//     QCoreApplication(int &argc, char **argv
// #ifndef Q_QDOC
//                      , int = ApplicationFlags
// #endif
//             );

//     ~QCoreApplication();

//     static QStringList arguments();

//     static void setAttribute(Qt::ApplicationAttribute attribute, bool on = true);
//     static bool testAttribute(Qt::ApplicationAttribute attribute);

//     static void setOrganizationDomain(const QString &orgDomain);
//     static QString organizationDomain();
//     static void setOrganizationName(const QString &orgName);
//     static QString organizationName();
//     static void setApplicationName(const QString &application);
//     static QString applicationName();
//     static void setApplicationVersion(const QString &version);
//     static QString applicationVersion();

//     static void setSetuidAllowed(bool allow);
//     static bool isSetuidAllowed();

//     static QCoreApplication *instance() { return self; }

// #ifndef QT_NO_QOBJECT
//     static int exec();
//     static void processEvents(QEventLoop::ProcessEventsFlags flags = QEventLoop::AllEvents);
//     static void processEvents(QEventLoop::ProcessEventsFlags flags, int maxtime);
//     static void exit(int retcode=0);

//     static bool sendEvent(QObject *receiver, QEvent *event);
//     static void postEvent(QObject *receiver, QEvent *event, int priority = Qt::NormalEventPriority);
//     static void sendPostedEvents(QObject *receiver = nullptr, int event_type = 0);
//     static void removePostedEvents(QObject *receiver, int eventType = 0);
// #if QT_DEPRECATED_SINCE(5, 3)
//     QT_DEPRECATED static bool hasPendingEvents();
// #endif
//     static QAbstractEventDispatcher *eventDispatcher();
//     static void setEventDispatcher(QAbstractEventDispatcher *eventDispatcher);

//     virtual bool notify(QObject *, QEvent *);

//     static bool startingUp();
//     static bool closingDown();
// #endif

//     static QString applicationDirPath();
//     static QString applicationFilePath();
//     static qint64 applicationPid() Q_DECL_CONST_FUNCTION;

// #if QT_CONFIG(library)
//     static void setLibraryPaths(const QStringList &);
//     static QStringList libraryPaths();
//     static void addLibraryPath(const QString &);
//     static void removeLibraryPath(const QString &);
// #endif // QT_CONFIG(library)

// #ifndef QT_NO_TRANSLATION
//     static bool installTranslator(QTranslator * messageFile);
//     static bool removeTranslator(QTranslator * messageFile);
// #endif

    static QString translate(const char * context,
                             const char * key,
                             const char * disambiguation = nullptr,
                             int n = -1);
// #if QT_DEPRECATED_SINCE(5, 0)
//     enum Encoding { UnicodeUTF8, Latin1, DefaultCodec = UnicodeUTF8, CodecForTr = UnicodeUTF8 };
//     QT_DEPRECATED static inline QString translate(const char * context, const char * key,
//                              const char * disambiguation, Encoding, int n = -1)
//         { return translate(context, key, disambiguation, n); }
// #endif

// #ifndef QT_NO_QOBJECT
// #  if QT_DEPRECATED_SINCE(5, 9)
//     QT_DEPRECATED static void flush();
// #  endif

//     void installNativeEventFilter(QAbstractNativeEventFilter *filterObj);
//     void removeNativeEventFilter(QAbstractNativeEventFilter *filterObj);

//     static bool isQuitLockEnabled();
//     static void setQuitLockEnabled(bool enabled);

// public Q_SLOTS:
//     static void quit();

// Q_SIGNALS:
//     void aboutToQuit(QPrivateSignal);

//     void organizationNameChanged();
//     void organizationDomainChanged();
//     void applicationNameChanged();
//     void applicationVersionChanged();

// protected:
//     bool event(QEvent *) override;

//     virtual bool compressEvent(QEvent *, QObject *receiver, QPostEventList *);
// #endif // QT_NO_QOBJECT

};


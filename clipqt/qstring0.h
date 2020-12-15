#include <QtCore>

class xQByteArray : public QString {
 public:
    char *data();
    inline int size() const;
};

class xQString : public QString {
 public:
    xQString(const char *ch);

    int length() const;
    QByteArray toUtf8() const;
};

class xQVariant : public QVariant {
public:
};

class xQSize : public QSize
{
 public:
    xQSize(int w, int h);

    inline int width() const noexcept;
    inline int height() const noexcept;

    inline int &rwidth() noexcept;
    inline int &rheight() noexcept;
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


class xQCoreApplication : public QCoreApplication {
public:
};


#include <QString>
#include <QChar>
#include <QSize>

class QString1  {
private:
    void* ptr;    
};

QString foo1() { return QString(); }
QChar foo2() { return QChar(); }
QCharRef foo3() {
    QString s("hhhh");
    return s[0];
}
QStringRef foo4() {
    QString s("hhhh");
    return s.leftRef(3);
}

QSize foo5() { return QSize(); }
QSizeF foo6() { return QSizeF(); }

double & foo7() { double val = 1.0; return val ; }

void foo8(const QString& a) {}

QString1 foo9() { return QString1(); }

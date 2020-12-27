#include <QString>
#include <QChar>
#include <QSize>
#include <QByteArray>
#include <QVariant>
#include <QJsonDocument>
#include <QJsonObject>
#include <QPoint>
#include <QPointF>
#include <QRect>
#include <QRectF>

class QString1  {
private:
    void* ptr;
};

QString foo1_QString() { return QString(); }
QChar foo2_QChar() { return QChar(); }
QCharRef foo3_QCharRef() {
    QString s("hhhh");
    return s[0];
}
QStringRef foo4_QStringRef() {
    QString s("hhhh");
    return s.leftRef(3);
}

QSize foo5_QSize() { return QSize(); }
QSizeF foo6_QSizeF() { return QSizeF(); }

double & foo7() { double val = 1.0; return val ; }

void foo8(const QString& a) {}

QString1 foo9_QString1() { return QString1(); }

QByteArray foo10_QByteArray() { return QByteArray(); }
QVariant foo11_QVariant() { return QVariant(); }
QJsonDocument foo12_QJsonDocument() { return QJsonDocument(); }
QJsonObject foo13_QJsonObject() { return QJsonObject(); }
QJsonValue foo14_QJsonValue() { return QJsonValue(); }

QPoint foo15_QPoint() { return QPoint(); }
QPointF foo16_QPointF() { return QPointF(); }

QRect foo17_QRect() {
    QRect* r1 = new QRect();
    delete r1;
    return QRect();
}
QRectF foo18_QRectF() { return QRectF(); }


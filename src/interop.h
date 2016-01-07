#ifndef _INTEROP_H_
#define _INTEROP_H_

#include <stdint.h>

#include <QList>
#include <QVector>
#include <QSet>
#include <QPair>
#include <QMap>
#include <QHash>
#include <QtCore>

template<class T>
char *irp_list_to_jsoin(QList<T> lst)
{
    QJsonObject jobj;
    QJsonArray jarr;
    for (T elem: &lst) {
        QJsonValue jv = QJsonValue::fromVariant(QVariant(elem));
        jarr.push_back(jv);
    }

    T avar = nullptr;
    int ety = QVariant(avar).type();

    jobj.insert("cname", "QList");
    jobj.insert("elems", jarr);
    jobj.insert("type", QJsonValue(ety));
    jobj.insert("tyname", QJsonValue(QVariant(avar).typeName()));

    QJsonDocument jdoc(jobj);
    char *res = strdup(jdoc.toJson().data());
    return res;
}

template<class T>
char *irp_vector_to_jsoin(QVector<T> lst)
{
    QJsonObject jobj;
    QJsonArray jarr;
    for (T elem: &lst) {
        QJsonValue jv = QJsonValue::fromVariant(QVariant(elem));
        jarr.push_back(jv);
    }

    T avar = nullptr;
    int ety = QVariant(avar).type();

    jobj.insert("cname", "QVector");
    jobj.insert("elems", jarr);
    jobj.insert("type", QJsonValue(ety));
    jobj.insert("tyname", QJsonValue(QVariant(avar).typeName()));

    QJsonDocument jdoc(jobj);
    char *res = strdup(jdoc.toJson().data());
    return res;
}

template<class T>
char *irp_set_to_jsoin(QSet<T> lst)
{
    QJsonObject jobj;
    QJsonArray jarr;
    for (T elem: &lst) {
        QJsonValue jv = QJsonValue::fromVariant(QVariant(elem));
        jarr.push_back(jv);
    }

    T avar = nullptr;
    int ety = QVariant(avar).type();

    jobj.insert("cname", "QSet");
    jobj.insert("elems", jarr);
    jobj.insert("type", QJsonValue(ety));
    jobj.insert("tyname", QJsonValue(QVariant(avar).typeName()));

    QJsonDocument jdoc(jobj);
    char *res = strdup(jdoc.toJson().data());
    return res;
}

template<class TK, class TV>
char *irp_map_to_jsoin(QMap<TK, TV> lst)
{
    QJsonObject jobj;
    QJsonObject jeobj;
    for (QPair<TK, TV> elem: &lst) {
        QJsonValue jk = QJsonValue::fromVariant(QVariant(elem.first));
        QJsonValue jv = QJsonValue::fromVariant(QVariant(elem.second));
        jeobj.insert(jk.toString(), jv);
    }

    TK kvar = nullptr;
    TV vvar = nullptr;
    int ekty = QVariant(kvar).type();
    int evty = QVariant(vvar).type();

    jobj.insert("cname", "QMap");
    jobj.insert("elems", jeobj);
    jobj.insert("ktype", QJsonValue(ekty));
    jobj.insert("ktyname", QJsonValue(QVariant(kvar).typeName()));
    jobj.insert("vtype", QJsonValue(evty));
    jobj.insert("vtyname", QJsonValue(QVariant(vvar).typeName()));

    QJsonDocument jdoc(jobj);
    char *res = strdup(jdoc.toJson().data());
    return res;
}

template<class TK, class TV>
    char *irp_hash_to_jsoin(QHash<TK, TV> lst)
{
    QJsonObject jobj;
    QJsonObject jeobj;
    for (QPair<TK, TV> elem: &lst) {
        QJsonValue jk = QJsonValue::fromVariant(QVariant(elem.first));
        QJsonValue jv = QJsonValue::fromVariant(QVariant(elem.second));
        jeobj.insert(jk.toString(), jv);
    }

    TK kvar = nullptr;
    TV vvar = nullptr;
    int ekty = QVariant(kvar).type();
    int evty = QVariant(vvar).type();

    jobj.insert("cname", "QHash");
    jobj.insert("elems", jeobj);
    jobj.insert("ktype", QJsonValue(ekty));
    jobj.insert("ktyname", QJsonValue(QVariant(kvar).typeName()));
    jobj.insert("vtype", QJsonValue(evty));
    jobj.insert("vtyname", QJsonValue(QVariant(vvar).typeName()));

    QJsonDocument jdoc(jobj);
    char *res = strdup(jdoc.toJson().data());
    return res;
}

template<class TK, class TV>
    char *irp_pair_to_jsoin(QPair<TK, TV> lst)
{
    QJsonObject jobj;
    QJsonObject jeobj;

    jeobj.insert(QJsonValue(QVariant(lst.first)).toString(),
                 QJsonValue(QVariant(lst.second)));

    TK kvar = nullptr;
    TV vvar = nullptr;
    int ekty = QVariant(kvar).type();
    int evty = QVariant(vvar).type();

    jobj.insert("cname", "QPair");
    jobj.insert("elems", jeobj);
    jobj.insert("ktype", QJsonValue(ekty));
    jobj.insert("ktyname", QJsonValue(QVariant(kvar).typeName()));
    jobj.insert("vtype", QJsonValue(evty));
    jobj.insert("vtyname", QJsonValue(QVariant(vvar).typeName()));

    QJsonDocument jdoc(jobj);
    char *res = strdup(jdoc.toJson().data());
    return res;
}


#endif

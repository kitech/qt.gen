#ifndef _QTHDRSRC_H_
#define _QTHDRSRC_H_

#include <string>  // fix: std::string 识别成了UNEXPOSED
// 使用在程序中自动生成的GEN_GO_QT_*宏确定是否#include某个模块，这样不需要修改这个文件了
// essentials
#ifdef GEN_GO_QT_CORE_LIB
#include <QtCore>
#endif
#ifdef GEN_GO_QT_GUI_LIB
#include <QtGui>
#endif
#ifdef GEN_GO_QT_WIDGETS_LIB
#include <QtWidgets>
#endif
#ifdef GEN_GO_QT_NETWORK_LIB
#include <QtNetwork>
#endif
#ifdef GEN_GO_QT_QML_LIB
#include <QtQml>
#endif
#ifdef GEN_GO_QT_QUICK_LIB
#include <QtQuick>
#endif

// add-ons
#ifdef GEN_GO_QT_QUICKTEMPLATES2_LIB
#include <QtQuickTemplates2>
#endif
#ifdef GEN_GO_QT_QUICKCONTROLS2_LIB
#include <QtQuickControls2>
#endif
#ifdef GEN_GO_QT_QUICKWIDGETS_LIB
#include <QtQuickWidgets>
#endif

// webengines
#ifdef GEN_GO_QT_WEBENGINECORE_LIB
#include <QtPositioning>
#endif
#ifdef GEN_GO_QT_WEBCHANNEL_LIB
#include <QtWebChannel>
#endif
#ifdef GEN_GO_QT_WEBENGINECORE_LIB
#include <QtWebEngineCore>
#endif
#ifdef GEN_GO_QT_WEBENGINE_LIB
#include <QtWebEngine>
#endif
#ifdef GEN_GO_QT_WEBENGINEWIDGETS_LIB
#include <QtWebEngineWidgets>
#endif

#ifdef GEN_GO_QT_SQL_LIB
#include <QtSql>
#endif
#ifdef GEN_GO_QT_MULTIMEDIA_LIB
#include <QtMultimedia>
#endif
#ifdef GEN_GO_QT_SVG_LIB
#include <QtSvg>
#endif
#ifdef GEN_GO_QT_MULTIMEDIAWIDGETS_LIB
#include <QtMultimediaWidgets>
#endif
#ifdef GEN_GO_QT_TEST_LIB
#include <QtTest>
#endif
#ifdef GEN_GO_QT_X11EXTRAS_LIB
#include <QtX11Extras>
#endif
#ifdef GEN_GO_QT_WINEXTRAS_LIB
#include <QtWinExtras>
#endif
#ifdef GEN_GO_QT_MACEXTRAS_LIB
#include <QtMacExtras>
#endif
#ifdef GEN_GO_QT_ANDROIDEXTRAS_LIB
#include <QtAndroidExtras>
#endif

// tools
#ifdef GEN_GO_QT_UITOOLS_LIB
#include <QtUiTools>
#endif

// #include <explicit_instantiate_templates.h>
// template class QFlags<int>;

typedef QList<QUrl> QUrlList;
typedef QList<QAbstractState*> QAbstractStateList;
typedef QList<QAccessibleInterface*> QAccessibleInterfaceList;
typedef QList<QSize> QSizeList;
typedef QList<QImageTextKeyLang> QImageTextKeyLangList;
typedef QList<QPolygonF> QPolygonFList;
typedef QList<QScreen *> QScreenList;
typedef QList<QStandardItem*> QStandardItemList;
typedef QList<QGlyphRun> QGlyphRunList;
typedef QList<QTextBlock> QTextBlockList;
typedef QList<QTextFrame *> QTextFrameList;
typedef QList<qreal>  qrealList;
typedef QList<QAction*> QActionList;
typedef QList<QKeySequence> QKeySequenceList;
typedef QList<QGraphicsWidget *> QGraphicsWidgetList;
typedef QList<QAbstractButton*> QAbstractButtonList;
typedef QList<int> intList;
typedef QList<QGesture *> QGestureList;
typedef QList<QGraphicsItem *> QGraphicsItemList;
typedef QList<QGraphicsTransform *> QGraphicsTransformList;
typedef QList<QListWidgetItem*> QListWidgetItemList;
typedef QList<QDockWidget*> QDockWidgetList;
typedef QList<QMdiSubWindow *> QMdiSubWindowList;
typedef QList<QScroller *> QScrollerList;
typedef QList<QTreeWidgetItem*> QTreeWidgetItemList;
typedef QList<QUndoStack*> QUndoStackList;
typedef QList<QNetworkConfiguration> QNetworkConfigurationList;
#ifdef QT_MULTIMEDIA_LIB
typedef QList<QCameraFocusZone> QCameraFocusZoneList;
typedef QList<QMediaResource> QMediaResourceList;
typedef QList<QCameraViewfinderSettings> QCameraViewfinderSettingsList;
typedef QList<QMediaContent> QMediaContentList;
#endif

/*
template class QHash<QString, QVariant>;
template class QMap<QString, QVariant>;
template class QHash<WId, QWidget *>;
template class QHash<int, QByteArray>;
template class QList<QVariant>;
template class QList<QByteArray>;
template class QList<QFileInfo>;
template class QList<QObject*>;
template class QList<QWindow*>;
template class QList<QWidget*>;
template class QList<QGraphicsItem *>;
template class QSet<QWidget*>;  // for QWidgetSet
template class QList<QModelIndex>;
template class QList<QQmlProperty>;
template class QList<QJSValue>;
template class QList<QQmlError>;
*/

// template class QList<QPointingDeviceUniqueId>;
// template class QHash<QString, QRemoteObjectSourceLocationInfo>;
// template class QList<QScriptValue>;
// template class QList<QDeclarativeProperty>;
// template class QList<QSurfaceDataRow*>;
// template class QList<QBarDataRow*>;
// template class QList<QCameraFocusZone>;
// template class QList<QMediaResource>;
// template class QMap<QModbusDataUnit::RegisterType, QModbusDataUnit>;

#include <QtCore/extra_export.h>

#endif

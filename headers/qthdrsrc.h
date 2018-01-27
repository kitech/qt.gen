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

#ifdef GEN_GO_QT_WEBENGINECORE_LIB
#include <QtWebEngineCore>
#endif
#ifdef GEN_GO_QT_WEBENGINE_LIB
#include <QtWebEngine>
#endif
#ifdef GEN_GO_QT_WEBCHANNEL_LIB
#include <QtWebChannel>
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


// template class QFlags<int>;
// #include <explicit_instantiate_templates.h>

#include <QtCore/extra_export.h>

#endif

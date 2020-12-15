#include <QtWidgets>

class xQWidget : public QWidget
{

public:
    explicit xQWidget(QWidget* parent = nullptr, Qt::WindowFlags f = Qt::WindowFlags());

    int devType() const override;

    WId winId() const;
    inline WId internalWinId() const;
    WId effectiveWinId() const;

/*     // GUI style setting */
/*     QStyle *style() const; */
/*     void setStyle(QStyle *); */
/*     // Widget types and states */

    bool isTopLevel() const;
    bool isWindow() const;

    bool isModal() const;
    Qt::WindowModality windowModality() const;
    void setWindowModality(Qt::WindowModality windowModality);

    bool isEnabled() const;
    bool isEnabledTo(const QWidget *) const;

public Q_SLOTS:
    void setEnabled(bool);
    void setDisabled(bool);
    void setWindowModified(bool);

    // Widget coordinates

public:
/*     QRect frameGeometry() const; */
/*     const QRect &geometry() const; */
/*     QRect normalGeometry() const; */

    int x() const;
    int y() const;
/*     QPoint pos() const; */
/*     QSize frameSize() const; */
/*     QSize size() const; */
    inline int width() const;
    inline int height() const;
/*     inline QRect rect() const; */
/*     QRect childrenRect() const; */
/*     QRegion childrenRegion() const; */

/*     QSize minimumSize() const; */
/*     QSize maximumSize() const; */
    int minimumWidth() const;
    int minimumHeight() const;
    int maximumWidth() const;
    int maximumHeight() const;
/*     void setMinimumSize(const QSize &); */
    void setMinimumSize(int minw, int minh);
/*     void setMaximumSize(const QSize &); */
    void setMaximumSize(int maxw, int maxh);
    void setMinimumWidth(int minw);
    void setMinimumHeight(int minh);
    void setMaximumWidth(int maxw);
    void setMaximumHeight(int maxh);

/*     QSize sizeIncrement() const; */
/*     void setSizeIncrement(const QSize &); */
/*     void setSizeIncrement(int w, int h); */
/*     QSize baseSize() const; */
/*     void setBaseSize(const QSize &); */
/*     void setBaseSize(int basew, int baseh); */

/*     void setFixedSize(const QSize &); */
    void setFixedSize(int w, int h);
    void setFixedWidth(int w);
    void setFixedHeight(int h);

/*     // Widget coordinate mapping */

/*     QPoint mapToGlobal(const QPoint &) const; */
/*     QPoint mapFromGlobal(const QPoint &) const; */
/*     QPoint mapToParent(const QPoint &) const; */
/*     QPoint mapFromParent(const QPoint &) const; */
/*     QPoint mapTo(const QWidget *, const QPoint &) const; */
/*     QPoint mapFrom(const QWidget *, const QPoint &) const; */

/*     QWidget *window() const; */
/*     QWidget *nativeParentWidget() const; */
/*     inline QWidget *topLevelWidget() const; */

/*     // Widget appearance functions */
/*     const QPalette &palette() const; */
/*     void setPalette(const QPalette &); */

/*     void setBackgroundRole(QPalette::ColorRole); */
/*     QPalette::ColorRole backgroundRole() const; */

/*     void setForegroundRole(QPalette::ColorRole); */
/*     QPalette::ColorRole foregroundRole() const; */

/*     const QFont &font() const; */
/*     void setFont(const QFont &); */
/*     QFontMetrics fontMetrics() const; */
/*     QFontInfo fontInfo() const; */

/* #ifndef QT_NO_CURSOR */
/*     QCursor cursor() const; */
/*     void setCursor(const QCursor &); */
/*     void unsetCursor(); */
/* #endif */

/*     void setMouseTracking(bool enable); */
/*     bool hasMouseTracking() const; */
/*     bool underMouse() const; */

/*     void setTabletTracking(bool enable); */
/*     bool hasTabletTracking() const; */

/*     void setMask(const QBitmap &); */
/*     void setMask(const QRegion &); */
/*     QRegion mask() const; */
/*     void clearMask(); */

/*     void render(QPaintDevice *target, const QPoint &targetOffset = QPoint(), */
/*                 const QRegion &sourceRegion = QRegion(), */
/*                 RenderFlags renderFlags = RenderFlags(DrawWindowBackground | DrawChildren)); */

/*     void render(QPainter *painter, const QPoint &targetOffset = QPoint(), */
/*                 const QRegion &sourceRegion = QRegion(), */
/*                 RenderFlags renderFlags = RenderFlags(DrawWindowBackground | DrawChildren)); */

/*     Q_INVOKABLE QPixmap grab(const QRect &rectangle = QRect(QPoint(0, 0), QSize(-1, -1))); */

/* #if QT_CONFIG(graphicseffect) */
/*     QGraphicsEffect *graphicsEffect() const; */
/*     void setGraphicsEffect(QGraphicsEffect *effect); */
/* #endif // QT_CONFIG(graphicseffect) */

/* #ifndef QT_NO_GESTURES */
/*     void grabGesture(Qt::GestureType type, Qt::GestureFlags flags = Qt::GestureFlags()); */
/*     void ungrabGesture(Qt::GestureType type); */
/* #endif */

public Q_SLOTS:
    void setWindowTitle(const QString &);
#ifndef QT_NO_STYLE_STYLESHEET
    void setStyleSheet(const QString& styleSheet);
#endif
public:
#ifndef QT_NO_STYLE_STYLESHEET
    QString styleSheet() const;
#endif
    QString windowTitle() const;
/*     void setWindowIcon(const QIcon &icon); */
/*     QIcon windowIcon() const; */
    void setWindowIconText(const QString &);
    QString windowIconText() const;
    void setWindowRole(const QString &);
    QString windowRole() const;
    void setWindowFilePath(const QString &filePath);
    QString windowFilePath() const;

    void setWindowOpacity(qreal level);
    qreal windowOpacity() const;

    bool isWindowModified() const;
#ifndef QT_NO_TOOLTIP
    void setToolTip(const QString &);
    QString toolTip() const;
    void setToolTipDuration(int msec);
    int toolTipDuration() const;
#endif
#if QT_CONFIG(statustip)
    void setStatusTip(const QString &);
    QString statusTip() const;
#endif
#if QT_CONFIG(whatsthis)
    void setWhatsThis(const QString &);
    QString whatsThis() const;
#endif

/*     void setLayoutDirection(Qt::LayoutDirection direction); */
/*     Qt::LayoutDirection layoutDirection() const; */
/*     void unsetLayoutDirection(); */

/*     void setLocale(const QLocale &locale); */
/*     QLocale locale() const; */
/*     void unsetLocale(); */

/*     inline bool isRightToLeft() const; */
/*     inline bool isLeftToRight() const; */

public Q_SLOTS:
    inline void setFocus();

public:
    bool isActiveWindow() const;
    void activateWindow();
    void clearFocus();

/*     void setFocus(Qt::FocusReason reason); */
/*     Qt::FocusPolicy focusPolicy() const; */
/*     void setFocusPolicy(Qt::FocusPolicy policy); */
/*     bool hasFocus() const; */
/*     static void setTabOrder(QWidget *, QWidget *); */
/*     void setFocusProxy(QWidget *); */
/*     QWidget *focusProxy() const; */
/*     Qt::ContextMenuPolicy contextMenuPolicy() const; */
/*     void setContextMenuPolicy(Qt::ContextMenuPolicy policy); */

/*     // Grab functions */
/*     void grabMouse(); */
/* #ifndef QT_NO_CURSOR */
/*     void grabMouse(const QCursor &); */
/* #endif */
/*     void releaseMouse(); */
/*     void grabKeyboard(); */
/*     void releaseKeyboard(); */
/* #ifndef QT_NO_SHORTCUT */
/*     int grabShortcut(const QKeySequence &key, Qt::ShortcutContext context = Qt::WindowShortcut); */
/*     void releaseShortcut(int id); */
/*     void setShortcutEnabled(int id, bool enable = true); */
/*     void setShortcutAutoRepeat(int id, bool enable = true); */
/* #endif */
/*     static QWidget *mouseGrabber(); */
/*     static QWidget *keyboardGrabber(); */

/*     // Update/refresh functions */
/*     inline bool updatesEnabled() const; */
/*     void setUpdatesEnabled(bool enable); */

/* #if QT_CONFIG(graphicsview) */
/*     QGraphicsProxyWidget *graphicsProxyWidget() const; */
/* #endif */

public Q_SLOTS:
    void update();
    void repaint();

/* public: */
/*     inline void update(int x, int y, int w, int h); */
/*     void update(const QRect&); */
/*     void update(const QRegion&); */

/*     void repaint(int x, int y, int w, int h); */
/*     void repaint(const QRect &); */
/*     void repaint(const QRegion &); */

public Q_SLOTS:
    // Widget management functions

    virtual void setVisible(bool visible);
    void setHidden(bool hidden);
    void show();
    void hide();

    void showMinimized();
    void showMaximized();
    void showFullScreen();
    void showNormal();

    bool close();
    void raise();
    void lower();

public:
/*     void stackUnder(QWidget*); */
/*     void move(int x, int y); */
/*     void move(const QPoint &); */
/*     void resize(int w, int h); */
/*     void resize(const QSize &); */
/*     inline void setGeometry(int x, int y, int w, int h); */
/*     void setGeometry(const QRect &); */
/*     QByteArray saveGeometry() const; */
/*     bool restoreGeometry(const QByteArray &geometry); */
    void adjustSize();
    bool isVisible() const;
    bool isVisibleTo(const QWidget *) const;
    inline bool isHidden() const;

    bool isMinimized() const;
    bool isMaximized() const;
    bool isFullScreen() const;

/*     Qt::WindowStates windowState() const; */
/*     void setWindowState(Qt::WindowStates state); */
/*     void overrideWindowState(Qt::WindowStates state); */

/*     virtual QSize sizeHint() const; */
/*     virtual QSize minimumSizeHint() const; */

/*     QSizePolicy sizePolicy() const; */
/*     void setSizePolicy(QSizePolicy); */
/*     inline void setSizePolicy(QSizePolicy::Policy horizontal, QSizePolicy::Policy vertical); */
/*     virtual int heightForWidth(int) const; */
/*     virtual bool hasHeightForWidth() const; */

/*     QRegion visibleRegion() const; */

/*     void setContentsMargins(int left, int top, int right, int bottom); */
/*     void setContentsMargins(const QMargins &margins); */
/*     QMargins contentsMargins() const; */

/*     QRect contentsRect() const; */

public:
    QLayout *layout() const;
    void setLayout(QLayout *);
    void updateGeometry();

/*     void setParent(QWidget *parent); */
/*     void setParent(QWidget *parent, Qt::WindowFlags f); */

/*     void scroll(int dx, int dy); */
/*     void scroll(int dx, int dy, const QRect&); */

/*     // Misc. functions */

/*     QWidget *focusWidget() const; */
/*     QWidget *nextInFocusChain() const; */
/*     QWidget *previousInFocusChain() const; */

/*     // drag and drop */
/*     bool acceptDrops() const; */
/*     void setAcceptDrops(bool on); */

/* #ifndef QT_NO_ACTION */
/*     //actions */
/*     void addAction(QAction *action); */
/* #if QT_VERSION >= QT_VERSION_CHECK(6,0,0) */
/*     void addActions(const QList<QAction*> &actions); */
/*     void insertActions(QAction *before, const QList<QAction*> &actions); */
/* #else */
/*     void addActions(QList<QAction*> actions); */
/*     void insertActions(QAction *before, QList<QAction*> actions); */
/* #endif */
/*     void insertAction(QAction *before, QAction *action); */
/*     void removeAction(QAction *action); */
/*     // QList<QAction*> actions() const; */
/* #endif */

    QWidget *parentWidget() const;

/*     void setWindowFlags(Qt::WindowFlags type); */
/*     inline Qt::WindowFlags windowFlags() const; */
/*     void setWindowFlag(Qt::WindowType, bool on = true); */
/*     void overrideWindowFlags(Qt::WindowFlags type); */

/*     inline Qt::WindowType windowType() const; */

/*     static QWidget *find(WId); */
/*     inline QWidget *childAt(int x, int y) const; */
/*     QWidget *childAt(const QPoint &p) const; */

/*     void setAttribute(Qt::WidgetAttribute, bool on = true); */
/*     inline bool testAttribute(Qt::WidgetAttribute) const; */

/*     QPaintEngine *paintEngine() const override; */

/*     void ensurePolished() const; */

/*     bool isAncestorOf(const QWidget *child) const; */

/* #ifdef QT_KEYPAD_NAVIGATION */
/*     bool hasEditFocus() const; */
/*     void setEditFocus(bool on); */
/* #endif */

/*     bool autoFillBackground() const; */
/*     void setAutoFillBackground(bool enabled); */

/*     QBackingStore *backingStore() const; */

/*     QWindow *windowHandle() const; */
/*     QScreen *screen() const; */

/*     static QWidget *createWindowContainer(QWindow *window, QWidget *parent=nullptr, Qt::WindowFlags flags=Qt::WindowFlags()); */

/* Q_SIGNALS: */
/*     void windowTitleChanged(const QString &title); */
/*     void windowIconChanged(const QIcon &icon); */
/*     void windowIconTextChanged(const QString &iconText); */
/*     void customContextMenuRequested(const QPoint &pos); */

/* protected: */
/*     // Event handlers */
/*     bool event(QEvent *event) override; */
/*     virtual void mousePressEvent(QMouseEvent *event); */
/*     virtual void mouseReleaseEvent(QMouseEvent *event); */
/*     virtual void mouseDoubleClickEvent(QMouseEvent *event); */
/*     virtual void mouseMoveEvent(QMouseEvent *event); */
/* #if QT_CONFIG(wheelevent) */
/*     virtual void wheelEvent(QWheelEvent *event); */
/* #endif */
/*     virtual void keyPressEvent(QKeyEvent *event); */
/*     virtual void keyReleaseEvent(QKeyEvent *event); */
/*     virtual void focusInEvent(QFocusEvent *event); */
/*     virtual void focusOutEvent(QFocusEvent *event); */
/*     virtual void enterEvent(QEvent *event); */
/*     virtual void leaveEvent(QEvent *event); */
/*     virtual void paintEvent(QPaintEvent *event); */
/*     virtual void moveEvent(QMoveEvent *event); */
/*     virtual void resizeEvent(QResizeEvent *event); */
/*     virtual void closeEvent(QCloseEvent *event); */
/* #ifndef QT_NO_CONTEXTMENU */
/*     virtual void contextMenuEvent(QContextMenuEvent *event); */
/* #endif */
/* #if QT_CONFIG(tabletevent) */
/*     virtual void tabletEvent(QTabletEvent *event); */
/* #endif */
/* #ifndef QT_NO_ACTION */
/*     virtual void actionEvent(QActionEvent *event); */
/* #endif */

/* #if QT_CONFIG(draganddrop) */
/*     virtual void dragEnterEvent(QDragEnterEvent *event); */
/*     virtual void dragMoveEvent(QDragMoveEvent *event); */
/*     virtual void dragLeaveEvent(QDragLeaveEvent *event); */
/*     virtual void dropEvent(QDropEvent *event); */
/* #endif */

/*     virtual void showEvent(QShowEvent *event); */
/*     virtual void hideEvent(QHideEvent *event); */

/* #if QT_VERSION >= QT_VERSION_CHECK(6, 0, 0) */
/*     virtual bool nativeEvent(const QByteArray &eventType, void *message, qintptr *result); */
/* #else */
/*     virtual bool nativeEvent(const QByteArray &eventType, void *message, long *result); */
/* #endif */

/*     // Misc. protected functions */
/*     virtual void changeEvent(QEvent *); */

};

class xQLayout : public QLayout {
 public:

}; // QLayout

class xQLayoutItem : public QLayoutItem {
 public:

}; // QLayoutItem

class xQSpacerItem : public QSpacerItem
{
 public:
    xQSpacerItem(int w, int h,
             QSizePolicy::Policy hData = QSizePolicy::Minimum,
             QSizePolicy::Policy vData = QSizePolicy::Minimum);
    //~QSpacerItem();

    void changeSize(int w, int h,
                    QSizePolicy::Policy hData = QSizePolicy::Minimum,
                    QSizePolicy::Policy vData = QSizePolicy::Minimum);
    QSize sizeHint() const override;
    QSize minimumSize() const override;
    QSize maximumSize() const override;
    Qt::Orientations expandingDirections() const override;
    bool isEmpty() const override;
    // void setGeometry(const QRect&) override;
    /* QRect geometry() const override; */
    /* QSpacerItem *spacerItem() override; */
    /* QSizePolicy sizePolicy() const { return sizeP; } */

}; // QSpacerItem

class xQWidgetItem : public QWidgetItem
{
public:
    explicit xQWidgetItem(QWidget *w) : wid(w) { }
    //~QWidgetItem();

    QSize sizeHint() const override;
    QSize minimumSize() const override;
    QSize maximumSize() const override;
    Qt::Orientations expandingDirections() const override;
    bool isEmpty() const override;
    //void setGeometry(const QRect&) override;
/*     QRect geometry() const override; */
/* #if QT_VERSION < QT_VERSION_CHECK(6, 0, 0) */
/*     QWidget *widget() override; */
/* #else */
/*     QWidget *widget() const override; */
/* #endif */

/*     bool hasHeightForWidth() const override; */
/*     int heightForWidth(int) const override; */
/*     QSizePolicy::ControlTypes controlTypes() const override; */
/* protected: */
/*     QWidget *wid; */
}; // QWidgetItem

class xQWidgetItemV2 : public QWidgetItemV2
{
public:
    explicit xQWidgetItemV2(QWidget *widget);
    //~QWidgetItemV2();

    QSize sizeHint() const override;
    QSize minimumSize() const override;
    QSize maximumSize() const override;
    int heightForWidth(int width) const override;

}; // QWidgetItem2


class xQBoxLayout : public QBoxLayout
{
 public:
    explicit xQBoxLayout(Direction, QWidget *parent = nullptr);

    //~QBoxLayout();

    Direction direction() const;
    void setDirection(Direction);

    void addSpacing(int size);
    void addStretch(int stretch = 0);
    /* void addSpacerItem(QSpacerItem *spacerItem); */
    /* void addWidget(QWidget *, int stretch = 0, Qt::Alignment alignment = Qt::Alignment()); */
    /* void addLayout(QLayout *layout, int stretch = 0); */
    /* void addStrut(int); */
    /* void addItem(QLayoutItem *) override; */

    /* void insertSpacing(int index, int size); */
    /* void insertStretch(int index, int stretch = 0); */
    /* void insertSpacerItem(int index, QSpacerItem *spacerItem); */
    /* void insertWidget(int index, QWidget *widget, int stretch = 0, Qt::Alignment alignment = Qt::Alignment()); */
    /* void insertLayout(int index, QLayout *layout, int stretch = 0); */
    /* void insertItem(int index, QLayoutItem *); */

    /* int spacing() const; */
    /* void setSpacing(int spacing); */

    /* bool setStretchFactor(QWidget *w, int stretch); */
    /* bool setStretchFactor(QLayout *l, int stretch); */
    /* void setStretch(int index, int stretch); */
    /* int stretch(int index) const; */

    /* QSize sizeHint() const override; */
    /* QSize minimumSize() const override; */
    /* QSize maximumSize() const override; */

    /* bool hasHeightForWidth() const override; */
    /* int heightForWidth(int) const override; */
    /* int minimumHeightForWidth(int) const override; */

    /* Qt::Orientations expandingDirections() const override; */
    /* void invalidate() override; */
    /* QLayoutItem *itemAt(int) const override; */
    /* QLayoutItem *takeAt(int) override; */
    /* int count() const override; */
    /* void setGeometry(const QRect&) override; */

}; // QBoxLayout

class xQHBoxLayout : public QHBoxLayout
{
public:
    xQHBoxLayout();
    explicit xQHBoxLayout(QWidget *parent);
    // ~QHBoxLayout();
};

class xQVBoxLayout : public QVBoxLayout
{
public:
    xQVBoxLayout();
    explicit xQVBoxLayout(QWidget *parent);
    // ~QVBoxLayout();
};


class xQButtonGroup : public QButtonGroup {
 public:
};

class xQAbstractButton : public QAbstractButton
{
 public:
    explicit xQAbstractButton(QWidget *parent = nullptr);

    void setText(const QString &text);
    QString text() const;

/*     void setIcon(const QIcon &icon); */
/*     QIcon icon() const; */

/*     QSize iconSize() const; */

/* #ifndef QT_NO_SHORTCUT */
/*     void setShortcut(const QKeySequence &key); */
/*     QKeySequence shortcut() const; */
/* #endif */

/*     void setCheckable(bool); */
/*     bool isCheckable() const; */

/*     bool isChecked() const; */

/*     void setDown(bool); */
/*     bool isDown() const; */

/*     void setAutoRepeat(bool); */
/*     bool autoRepeat() const; */

/*     void setAutoRepeatDelay(int); */
/*     int autoRepeatDelay() const; */

/*     void setAutoRepeatInterval(int); */
/*     int autoRepeatInterval() const; */

/*     void setAutoExclusive(bool); */
/*     bool autoExclusive() const; */

/* #if QT_CONFIG(buttongroup) */
/*     QButtonGroup *group() const; */
/* #endif */

/* public Q_SLOTS: */
/*     void setIconSize(const QSize &size); */
/*     void animateClick(int msec = 100); */
/*     void click(); */
/*     void toggle(); */
/*     void setChecked(bool); */

/* Q_SIGNALS: */
/*     void pressed(); */
/*     void released(); */
/*     void clicked(bool checked = false); */
/*     void toggled(bool checked); */

/* protected: */
/*     void paintEvent(QPaintEvent *e) override = 0; */
/*     virtual bool hitButton(const QPoint &pos) const; */
/*     virtual void checkStateSet(); */
/*     virtual void nextCheckState(); */

/*     bool event(QEvent *e) override; */
/*     void keyPressEvent(QKeyEvent *e) override; */
/*     void keyReleaseEvent(QKeyEvent *e) override; */
/*     void mousePressEvent(QMouseEvent *e) override; */
/*     void mouseReleaseEvent(QMouseEvent *e) override; */
/*     void mouseMoveEvent(QMouseEvent *e) override; */
/*     void focusInEvent(QFocusEvent *e) override; */
/*     void focusOutEvent(QFocusEvent *e) override; */
/*     void changeEvent(QEvent *e) override; */
/*     void timerEvent(QTimerEvent *e) override; */

};

class xQPushButton : public QPushButton {
 public:
    explicit xQPushButton(QWidget *parent = nullptr);
    explicit xQPushButton(const QString &text, QWidget *parent = nullptr);

    QSize sizeHint() const override;
    QSize minimumSizeHint() const override;

    bool autoDefault() const;
    void setAutoDefault(bool);
    bool isDefault() const;
    void setDefault(bool);

/* #if QT_CONFIG(menu) */
/*     void setMenu(QMenu* menu); */
/*     QMenu* menu() const; */
/* #endif */

/*     void setFlat(bool); */
/*     bool isFlat() const; */

/*     public Q_SLOTS: */
/* #if QT_CONFIG(menu) */
/*         void showMenu(); */
/* #endif */

/*  protected: */
/*         bool event(QEvent *e) override; */
/*         void paintEvent(QPaintEvent *) override; */
/*         void keyPressEvent(QKeyEvent *) override; */
/*         void focusInEvent(QFocusEvent *) override; */
/*         void focusOutEvent(QFocusEvent *) override; */
/*         void initStyleOption(QStyleOptionButton *option) const; */
/*         bool hitButton(const QPoint &pos) const override; */
/*         QPushButton(QPushButtonPrivate &dd, QWidget* parent = nullptr); */

}; // QPushButton

class  xQCheckBox : public QCheckBox
{
public:
    explicit xQCheckBox(QWidget *parent = nullptr);
    explicit xQCheckBox(const QString &text, QWidget *parent = nullptr);
    //    ~QCheckBox();

    QSize sizeHint() const override;
    QSize minimumSizeHint() const override;

    void setTristate(bool y = true);
    bool isTristate() const;

    Qt::CheckState checkState() const;
    void setCheckState(Qt::CheckState state);

Q_SIGNALS:
    void stateChanged(int);

/* protected: */
/*     bool event(QEvent *e) override; */
/*     bool hitButton(const QPoint &pos) const override; */
/*     void checkStateSet() override; */
/*     void nextCheckState() override; */
/*     void paintEvent(QPaintEvent *) override; */
/*     void mouseMoveEvent(QMouseEvent *) override; */
/*     void initStyleOption(QStyleOptionButton *option) const; */


}; // QCheckBox

class xQRadioButton : public QRadioButton
{
 public:
    explicit xQRadioButton(QWidget *parent = nullptr);
    explicit xQRadioButton(const QString &text, QWidget *parent = nullptr);
    //    ~QRadioButton();

    QSize sizeHint() const override;
    QSize minimumSizeHint() const override;

 /* protected: */
 /*    bool event(QEvent *e) override; */
 /*    bool hitButton(const QPoint &) const override; */
 /*    void paintEvent(QPaintEvent *) override; */
 /*    void mouseMoveEvent(QMouseEvent *) override; */
 /*    void initStyleOption(QStyleOptionButton *button) const; */

}; // QRadioButton


class xQFrame : public QFrame
{
public:
    explicit xQFrame(QWidget* parent = nullptr, Qt::WindowFlags f = Qt::WindowFlags());
    // ~QFrame();

    int frameStyle() const;
    void setFrameStyle(int);

    int frameWidth() const;

    QSize sizeHint() const override;

/*     enum Shape { */
/*         NoFrame  = 0, // no frame */
/*         Box = 0x0001, // rectangular box */
/*         Panel = 0x0002, // rectangular panel */
/*         WinPanel = 0x0003, // rectangular panel (Windows) */
/*         HLine = 0x0004, // horizontal line */
/*         VLine = 0x0005, // vertical line */
/*         StyledPanel = 0x0006 // rectangular panel depending on the GUI style */
/*     }; */
/*     Q_ENUM(Shape) */
/*     enum Shadow { */
/*         Plain = 0x0010, // plain line */
/*         Raised = 0x0020, // raised shadow effect */
/*         Sunken = 0x0030 // sunken shadow effect */
/*     }; */
/*     Q_ENUM(Shadow) */

/*     enum StyleMask { */
/*         Shadow_Mask = 0x00f0, // mask for the shadow */
/*         Shape_Mask = 0x000f // mask for the shape */
/*     }; */

/*     Shape frameShape() const; */
/*     void setFrameShape(Shape); */
/*     Shadow frameShadow() const; */
/*     void setFrameShadow(Shadow); */

/*     int lineWidth() const; */
/*     void setLineWidth(int); */

/*     int midLineWidth() const; */
/*     void setMidLineWidth(int); */

/*     QRect frameRect() const; */
/*     void setFrameRect(const QRect &); */

/* protected: */
/*     bool event(QEvent *e) override; */
/*     void paintEvent(QPaintEvent *) override; */
/*     void changeEvent(QEvent *) override; */
/*     void drawFrame(QPainter *); */


/* protected: */
/*     QFrame(QFramePrivate &dd, QWidget* parent = nullptr, Qt::WindowFlags f = Qt::WindowFlags()); */
/*     void initStyleOption(QStyleOptionFrame *option) const; */

}; // QFrame


class xQLabel : public QLabel
{

public:
    explicit xQLabel(QWidget *parent=nullptr, Qt::WindowFlags f=Qt::WindowFlags());
    explicit xQLabel(const QString &text, QWidget *parent=nullptr, Qt::WindowFlags f=Qt::WindowFlags());
    //~QLabel();

    QString text() const;

/* #if QT_DEPRECATED_SINCE(5,15) */
/*     QT_DEPRECATED_VERSION_X(5, 15, "Use the other overload which returns QPixmap by-value") */
/*     const QPixmap *pixmap() const; // ### Qt 7: Remove function */

/*     QPixmap pixmap(Qt::ReturnByValueConstant) const; */
/* #else */
/*     QPixmap pixmap(Qt::ReturnByValueConstant = Qt::ReturnByValue) const; // ### Qt 7: Remove arg */
/* #endif // QT_DEPRECATED_SINCE(5,15) */

/* #ifndef QT_NO_PICTURE */
/* #  if QT_DEPRECATED_SINCE(5,15) */
/*     QT_DEPRECATED_VERSION_X(5, 15, "Use the other overload which returns QPicture by-value") */
/*     const QPicture *picture() const; // ### Qt 7: Remove function */

/*     QPicture picture(Qt::ReturnByValueConstant) const; */
/* #  else */
/*     QPicture picture(Qt::ReturnByValueConstant = Qt::ReturnByValue) const; // ### Qt 7: Remove arg */
/* #  endif // QT_DEPRECATED_SINCE(5,15) */
/* #endif */
/* #if QT_CONFIG(movie) */
/*     QMovie *movie() const; */
/* #endif */

/*     Qt::TextFormat textFormat() const; */
/*     void setTextFormat(Qt::TextFormat); */

/*     Qt::Alignment alignment() const; */
/*     void setAlignment(Qt::Alignment); */

/*     void setWordWrap(bool on); */
/*     bool wordWrap() const; */

/*     int indent() const; */
/*     void setIndent(int); */

/*     int margin() const; */
/*     void setMargin(int); */

/*     bool hasScaledContents() const; */
/*     void setScaledContents(bool); */
/*     QSize sizeHint() const override; */
/*     QSize minimumSizeHint() const override; */
/* #ifndef QT_NO_SHORTCUT */
/*     void setBuddy(QWidget *); */
/*     QWidget *buddy() const; */
/* #endif */
/*     int heightForWidth(int) const override; */

/*     bool openExternalLinks() const; */
/*     void setOpenExternalLinks(bool open); */

/*     void setTextInteractionFlags(Qt::TextInteractionFlags flags); */
/*     Qt::TextInteractionFlags textInteractionFlags() const; */

/*     void setSelection(int, int); */
/*     bool hasSelectedText() const; */
/*     QString selectedText() const; */
/*     int selectionStart() const; */

/* public Q_SLOTS: */
/*     void setText(const QString &); */
/*     void setPixmap(const QPixmap &); */
/* #ifndef QT_NO_PICTURE */
/*     void setPicture(const QPicture &); */
/* #endif */
/* #if QT_CONFIG(movie) */
/*     void setMovie(QMovie *movie); */
/* #endif */
/*     void setNum(int); */
/*     void setNum(double); */
/*     void clear(); */

/* Q_SIGNALS: */
/*     void linkActivated(const QString& link); */
/*     void linkHovered(const QString& link); */

/* protected: */
/*     bool event(QEvent *e) override; */
/*     void keyPressEvent(QKeyEvent *ev) override; */
/*     void paintEvent(QPaintEvent *) override; */
/*     void changeEvent(QEvent *) override; */
/*     void mousePressEvent(QMouseEvent *ev) override; */
/*     void mouseMoveEvent(QMouseEvent *ev) override; */
/*     void mouseReleaseEvent(QMouseEvent *ev) override; */
/* #ifndef QT_NO_CONTEXTMENU */
/*     void contextMenuEvent(QContextMenuEvent *ev) override; */
/* #endif // QT_NO_CONTEXTMENU */
/*     void focusInEvent(QFocusEvent *ev) override; */
/*     void focusOutEvent(QFocusEvent *ev) override; */
/*     bool focusNextPrevChild(bool next) override; */

}; // QLabel

class QLineEdit : public QWidget
{
 public:
    explicit xQLineEdit(QWidget *parent = nullptr);
    explicit xQLineEdit(const QString &, QWidget *parent = nullptr);
    //~QLineEdit();

    QString text() const;

    QString displayText() const;

    QString placeholderText() const;
    void setPlaceholderText(const QString &);

    int maxLength() const;
    void setMaxLength(int);

    void setFrame(bool);
    bool hasFrame() const;

    void setClearButtonEnabled(bool enable);
    bool isClearButtonEnabled() const;

/*     enum EchoMode { Normal, NoEcho, Password, PasswordEchoOnEdit }; */
/*     Q_ENUM(EchoMode) */
/*     EchoMode echoMode() const; */
/*     void setEchoMode(EchoMode); */

/*     bool isReadOnly() const; */
/*     void setReadOnly(bool); */

/* #ifndef QT_NO_VALIDATOR */
/*     void setValidator(const QValidator *); */
/*     const QValidator * validator() const; */
/* #endif */

/* #if QT_CONFIG(completer) */
/*     void setCompleter(QCompleter *completer); */
/*     QCompleter *completer() const; */
/* #endif */

/*     QSize sizeHint() const override; */
/*     QSize minimumSizeHint() const override; */

/*     int cursorPosition() const; */
/*     void setCursorPosition(int); */
/*     int cursorPositionAt(const QPoint &pos); */

/*     void setAlignment(Qt::Alignment flag); */
/*     Qt::Alignment alignment() const; */

/*     void cursorForward(bool mark, int steps = 1); */
/*     void cursorBackward(bool mark, int steps = 1); */
/*     void cursorWordForward(bool mark); */
/*     void cursorWordBackward(bool mark); */
/*     void backspace(); */
/*     void del(); */
/*     void home(bool mark); */
/*     void end(bool mark); */

/*     bool isModified() const; */
/*     void setModified(bool); */

/*     void setSelection(int, int); */
/*     bool hasSelectedText() const; */
/*     QString selectedText() const; */
/*     int selectionStart() const; */
/*     int selectionEnd() const; */
/*     int selectionLength() const; */

/*     bool isUndoAvailable() const; */
/*     bool isRedoAvailable() const; */

/*     void setDragEnabled(bool b); */
/*     bool dragEnabled() const; */

/*     void setCursorMoveStyle(Qt::CursorMoveStyle style); */
/*     Qt::CursorMoveStyle cursorMoveStyle() const; */

/*     QString inputMask() const; */
/*     void setInputMask(const QString &inputMask); */
/*     bool hasAcceptableInput() const; */

/*     void setTextMargins(int left, int top, int right, int bottom); */
/*     void setTextMargins(const QMargins &margins); */
/* #if QT_DEPRECATED_SINCE(5, 14) */
/*     QT_DEPRECATED_X("use textMargins()") */
/*     void getTextMargins(int *left, int *top, int *right, int *bottom) const; */
/* #endif */
/*     QMargins textMargins() const; */

/* #if QT_CONFIG(action) */
/*     using QWidget::addAction; */
/*     void addAction(QAction *action, ActionPosition position); */
/*     QAction *addAction(const QIcon &icon, ActionPosition position); */
/* #endif */

/* public Q_SLOTS: */
/*     void setText(const QString &); */
/*     void clear(); */
/*     void selectAll(); */
/*     void undo(); */
/*     void redo(); */
/* #ifndef QT_NO_CLIPBOARD */
/*     void cut(); */
/*     void copy() const; */
/*     void paste(); */
/* #endif */

/* public: */
/*     void deselect(); */
/*     void insert(const QString &); */
/* #ifndef QT_NO_CONTEXTMENU */
/*     QMenu *createStandardContextMenu(); */
/* #endif */

/* Q_SIGNALS: */
/*     void textChanged(const QString &); */
/*     void textEdited(const QString &); */
/*     void cursorPositionChanged(int, int); */
/*     void returnPressed(); */
/*     void editingFinished(); */
/*     void selectionChanged(); */
/*     void inputRejected(); */

/* protected: */
/*     void mousePressEvent(QMouseEvent *) override; */
/*     void mouseMoveEvent(QMouseEvent *) override; */
/*     void mouseReleaseEvent(QMouseEvent *) override; */
/*     void mouseDoubleClickEvent(QMouseEvent *) override; */
/*     void keyPressEvent(QKeyEvent *) override; */
/*     void focusInEvent(QFocusEvent *) override; */
/*     void focusOutEvent(QFocusEvent *) override; */
/*     void paintEvent(QPaintEvent *) override; */
/* #if QT_CONFIG(draganddrop) */
/*     void dragEnterEvent(QDragEnterEvent *) override; */
/*     void dragMoveEvent(QDragMoveEvent *e) override; */
/*     void dragLeaveEvent(QDragLeaveEvent *e) override; */
/*     void dropEvent(QDropEvent *) override; */
/* #endif */
/*     void changeEvent(QEvent *) override; */
/* #ifndef QT_NO_CONTEXTMENU */
/*     void contextMenuEvent(QContextMenuEvent *) override; */
/* #endif */

/*     void inputMethodEvent(QInputMethodEvent *) override; */
/*     void initStyleOption(QStyleOptionFrame *option) const; */
/* public: */
/*     QVariant inputMethodQuery(Qt::InputMethodQuery) const override; */
/*     Q_INVOKABLE QVariant inputMethodQuery(Qt::InputMethodQuery property, QVariant argument) const; */
/*     bool event(QEvent *) override; */
/* protected: */
/*     QRect cursorRect() const; */

public:

}; // end QLineEdit

class xQAbstractScrollArea : public QAbstractScrollArea
{
public:
    explicit xQAbstractScrollArea(QWidget *parent = nullptr);
    //    ~QAbstractScrollArea();

/*     enum SizeAdjustPolicy { */
/*         AdjustIgnored, */
/*         AdjustToContentsOnFirstShow, */
/*         AdjustToContents */
/*     }; */
/*     Q_ENUM(SizeAdjustPolicy) */

/*     Qt::ScrollBarPolicy verticalScrollBarPolicy() const; */
/*     void setVerticalScrollBarPolicy(Qt::ScrollBarPolicy); */
/*     QScrollBar *verticalScrollBar() const; */
/*     void setVerticalScrollBar(QScrollBar *scrollbar); */

/*     Qt::ScrollBarPolicy horizontalScrollBarPolicy() const; */
/*     void setHorizontalScrollBarPolicy(Qt::ScrollBarPolicy); */
/*     QScrollBar *horizontalScrollBar() const; */
/*     void setHorizontalScrollBar(QScrollBar *scrollbar); */

/*     QWidget *cornerWidget() const; */
/*     void setCornerWidget(QWidget *widget); */

/*     void addScrollBarWidget(QWidget *widget, Qt::Alignment alignment); */
/*     QWidgetList scrollBarWidgets(Qt::Alignment alignment); */

/*     QWidget *viewport() const; */
/*     void setViewport(QWidget *widget); */
/*     QSize maximumViewportSize() const; */

/*     QSize minimumSizeHint() const override; */

/*     QSize sizeHint() const override; */

/*     virtual void setupViewport(QWidget *viewport); */

/*     SizeAdjustPolicy sizeAdjustPolicy() const; */
/*     void setSizeAdjustPolicy(SizeAdjustPolicy policy); */

/* protected: */
/*     QAbstractScrollArea(QAbstractScrollAreaPrivate &dd, QWidget *parent = nullptr); */
/*     void setViewportMargins(int left, int top, int right, int bottom); */
/*     void setViewportMargins(const QMargins &margins); */
/*     QMargins viewportMargins() const; */

/*     bool eventFilter(QObject *, QEvent *) override; */
/*     bool event(QEvent *) override; */
/*     virtual bool viewportEvent(QEvent *); */

/*     void resizeEvent(QResizeEvent *) override; */
/*     void paintEvent(QPaintEvent *) override; */
/*     void mousePressEvent(QMouseEvent *) override; */
/*     void mouseReleaseEvent(QMouseEvent *) override; */
/*     void mouseDoubleClickEvent(QMouseEvent *) override; */
/*     void mouseMoveEvent(QMouseEvent *) override; */
/* #if QT_CONFIG(wheelevent) */
/*     void wheelEvent(QWheelEvent *) override; */
/* #endif */
/* #ifndef QT_NO_CONTEXTMENU */
/*     void contextMenuEvent(QContextMenuEvent *) override; */
/* #endif */
/* #if QT_CONFIG(draganddrop) */
/*     void dragEnterEvent(QDragEnterEvent *) override; */
/*     void dragMoveEvent(QDragMoveEvent *) override; */
/*     void dragLeaveEvent(QDragLeaveEvent *) override; */
/*     void dropEvent(QDropEvent *) override; */
/* #endif */

/*     void keyPressEvent(QKeyEvent *) override; */

/*     virtual void scrollContentsBy(int dx, int dy); */

/*     virtual QSize viewportSizeHint() const; */

}; // QAbstractScrollArea

class xQPlainTextEdit : public QPlainTextEdit
{
 public:
    explicit xQPlainTextEdit(QWidget *parent = nullptr);
    explicit xQPlainTextEdit(const QString &text, QWidget *parent = nullptr);
    //virtual ~QPlainTextEdit();

/*     void setDocument(QTextDocument *document); */
/*     QTextDocument *document() const; */

/*     void setPlaceholderText(const QString &placeholderText); */
/*     QString placeholderText() const; */

/*     void setTextCursor(const QTextCursor &cursor); */
/*     QTextCursor textCursor() const; */

/*     bool isReadOnly() const; */
/*     void setReadOnly(bool ro); */

/*     void setTextInteractionFlags(Qt::TextInteractionFlags flags); */
/*     Qt::TextInteractionFlags textInteractionFlags() const; */

/*     void mergeCurrentCharFormat(const QTextCharFormat &modifier); */
/*     void setCurrentCharFormat(const QTextCharFormat &format); */
/*     QTextCharFormat currentCharFormat() const; */

/*     bool tabChangesFocus() const; */
/*     void setTabChangesFocus(bool b); */

/*     inline void setDocumentTitle(const QString &title) */
/*     { document()->setMetaInformation(QTextDocument::DocumentTitle, title); } */
/*     inline QString documentTitle() const */
/*     { return document()->metaInformation(QTextDocument::DocumentTitle); } */

/*     inline bool isUndoRedoEnabled() const */
/*     { return document()->isUndoRedoEnabled(); } */
/*     inline void setUndoRedoEnabled(bool enable) */
/*     { document()->setUndoRedoEnabled(enable); } */

/*     inline void setMaximumBlockCount(int maximum) */
/*     { document()->setMaximumBlockCount(maximum); } */
/*     inline int maximumBlockCount() const */
/*     { return document()->maximumBlockCount(); } */


/*     LineWrapMode lineWrapMode() const; */
/*     void setLineWrapMode(LineWrapMode mode); */

/*     QTextOption::WrapMode wordWrapMode() const; */
/*     void setWordWrapMode(QTextOption::WrapMode policy); */

/*     void setBackgroundVisible(bool visible); */
/*     bool backgroundVisible() const; */

/*     void setCenterOnScroll(bool enabled); */
/*     bool centerOnScroll() const; */

/*     bool find(const QString &exp, QTextDocument::FindFlags options = QTextDocument::FindFlags()); */
/* #ifndef QT_NO_REGEXP */
/*     bool find(const QRegExp &exp, QTextDocument::FindFlags options = QTextDocument::FindFlags()); */
/* #endif */
/* #if QT_CONFIG(regularexpression) */
/*     bool find(const QRegularExpression &exp, QTextDocument::FindFlags options = QTextDocument::FindFlags()); */
/* #endif */

/*     inline QString toPlainText() const */
/*     { return document()->toPlainText(); } */

/*     void ensureCursorVisible(); */

/*     virtual QVariant loadResource(int type, const QUrl &name); */
/* #ifndef QT_NO_CONTEXTMENU */
/*     QMenu *createStandardContextMenu(); */
/*     QMenu *createStandardContextMenu(const QPoint &position); */
/* #endif */

/*     QTextCursor cursorForPosition(const QPoint &pos) const; */
/*     QRect cursorRect(const QTextCursor &cursor) const; */
/*     QRect cursorRect() const; */

/*     QString anchorAt(const QPoint &pos) const; */

/*     bool overwriteMode() const; */
/*     void setOverwriteMode(bool overwrite); */

/* #if QT_DEPRECATED_SINCE(5, 10) */
/*     QT_DEPRECATED int tabStopWidth() const; */
/*     QT_DEPRECATED void setTabStopWidth(int width); */
/* #endif */

/*     qreal tabStopDistance() const; */
/*     void setTabStopDistance(qreal distance); */

/*     int cursorWidth() const; */
/*     void setCursorWidth(int width); */

/*     void setExtraSelections(const QList<QTextEdit::ExtraSelection> &selections); */
/*     QList<QTextEdit::ExtraSelection> extraSelections() const; */

/*     void moveCursor(QTextCursor::MoveOperation operation, QTextCursor::MoveMode mode = QTextCursor::MoveAnchor); */

/*     bool canPaste() const; */

/*     void print(QPagedPaintDevice *printer) const; */

/*     int blockCount() const; */
/*     QVariant inputMethodQuery(Qt::InputMethodQuery property) const override; */
/*     Q_INVOKABLE QVariant inputMethodQuery(Qt::InputMethodQuery query, QVariant argument) const; */

/* public Q_SLOTS: */

/*     void setPlainText(const QString &text); */

/* #ifndef QT_NO_CLIPBOARD */
/*     void cut(); */
/*     void copy(); */
/*     void paste(); */
/* #endif */

/*     void undo(); */
/*     void redo(); */

/*     void clear(); */
/*     void selectAll(); */

/*     void insertPlainText(const QString &text); */

/*     void appendPlainText(const QString &text); */
/*     void appendHtml(const QString &html); */

/*     void centerCursor(); */

/*     void zoomIn(int range = 1); */
/*     void zoomOut(int range = 1); */

/* Q_SIGNALS: */
/*     void textChanged(); */
/*     void undoAvailable(bool b); */
/*     void redoAvailable(bool b); */
/*     void copyAvailable(bool b); */
/*     void selectionChanged(); */
/*     void cursorPositionChanged(); */

/*     void updateRequest(const QRect &rect, int dy); */
/*     void blockCountChanged(int newBlockCount); */
/*     void modificationChanged(bool); */

/* protected: */
/*     virtual bool event(QEvent *e) override; */
/*     virtual void timerEvent(QTimerEvent *e) override; */
/*     virtual void keyPressEvent(QKeyEvent *e) override; */
/*     virtual void keyReleaseEvent(QKeyEvent *e) override; */
/*     virtual void resizeEvent(QResizeEvent *e) override; */
/*     virtual void paintEvent(QPaintEvent *e) override; */
/*     virtual void mousePressEvent(QMouseEvent *e) override; */
/*     virtual void mouseMoveEvent(QMouseEvent *e) override; */
/*     virtual void mouseReleaseEvent(QMouseEvent *e) override; */
/*     virtual void mouseDoubleClickEvent(QMouseEvent *e) override; */
/*     virtual bool focusNextPrevChild(bool next) override; */
/* #ifndef QT_NO_CONTEXTMENU */
/*     virtual void contextMenuEvent(QContextMenuEvent *e) override; */
/* #endif */
/* #if QT_CONFIG(draganddrop) */
/*     virtual void dragEnterEvent(QDragEnterEvent *e) override; */
/*     virtual void dragLeaveEvent(QDragLeaveEvent *e) override; */
/*     virtual void dragMoveEvent(QDragMoveEvent *e) override; */
/*     virtual void dropEvent(QDropEvent *e) override; */
/* #endif */
/*     virtual void focusInEvent(QFocusEvent *e) override; */
/*     virtual void focusOutEvent(QFocusEvent *e) override; */
/*     virtual void showEvent(QShowEvent *) override; */
/*     virtual void changeEvent(QEvent *e) override; */
/* #if QT_CONFIG(wheelevent) */
/*     virtual void wheelEvent(QWheelEvent *e) override; */
/* #endif */

/*     virtual QMimeData *createMimeDataFromSelection() const; */
/*     virtual bool canInsertFromMimeData(const QMimeData *source) const; */
/*     virtual void insertFromMimeData(const QMimeData *source); */

/*     virtual void inputMethodEvent(QInputMethodEvent *) override; */

/*     QPlainTextEdit(QPlainTextEditPrivate &dd, QWidget *parent); */

/*     virtual void scrollContentsBy(int dx, int dy) override; */
/*     virtual void doSetTextCursor(const QTextCursor &cursor); */

/*     QTextBlock firstVisibleBlock() const; */
/*     QPointF contentOffset() const; */
/*     QRectF blockBoundingRect(const QTextBlock &block) const; */
/*     QRectF blockBoundingGeometry(const QTextBlock &block) const; */
/*     QAbstractTextDocumentLayout::PaintContext getPaintContext() const; */

/*     void zoomInF(float range); */

}; // QPlainTextEdit

class xQComboBox : public QComboBox
{
public:
    explicit xQComboBox(QWidget *parent = nullptr);
    //    ~QComboBox();

    int maxVisibleItems() const;
    void setMaxVisibleItems(int maxItems);

    int count() const;
    void setMaxCount(int max);
    int maxCount() const;

/* #if QT_CONFIG(completer) */
/* #if QT_DEPRECATED_SINCE(5, 13) */
/*     QT_DEPRECATED_X("Use completer() instead.") */
/*     bool autoCompletion() const; */
/*     QT_DEPRECATED_X("Use setCompleter() instead.") */
/*     void setAutoCompletion(bool enable); */
/*     QT_DEPRECATED_X("Use completer()->caseSensitivity() instead.") */
/*     Qt::CaseSensitivity autoCompletionCaseSensitivity() const; */
/*     QT_DEPRECATED_X("Use completer()->setCaseSensitivity() instead.") */
/*     void setAutoCompletionCaseSensitivity(Qt::CaseSensitivity sensitivity); */
/* #endif */
/* #endif */

/*     bool duplicatesEnabled() const; */
/*     void setDuplicatesEnabled(bool enable); */

/*     void setFrame(bool); */
/*     bool hasFrame() const; */

/*     inline int findText(const QString &text, */
/*                         Qt::MatchFlags flags = static_cast<Qt::MatchFlags>(Qt::MatchExactly|Qt::MatchCaseSensitive)) const */
/*         { return findData(text, Qt::DisplayRole, flags); } */
/*     int findData(const QVariant &data, int role = Qt::UserRole, */
/*                  Qt::MatchFlags flags = static_cast<Qt::MatchFlags>(Qt::MatchExactly|Qt::MatchCaseSensitive)) const; */

/*     enum InsertPolicy { */
/*         NoInsert, */
/*         InsertAtTop, */
/*         InsertAtCurrent, */
/*         InsertAtBottom, */
/*         InsertAfterCurrent, */
/*         InsertBeforeCurrent, */
/*         InsertAlphabetically */
/*     }; */
/*     Q_ENUM(InsertPolicy) */

/*     InsertPolicy insertPolicy() const; */
/*     void setInsertPolicy(InsertPolicy policy); */

/*     enum SizeAdjustPolicy { */
/*         AdjustToContents, */
/*         AdjustToContentsOnFirstShow, */
/* #if QT_DEPRECATED_SINCE(5, 15) */
/*         AdjustToMinimumContentsLength Q_DECL_ENUMERATOR_DEPRECATED_X( */
/*             "Use AdjustToContents or AdjustToContentsOnFirstShow"), // ### Qt 6: remove */
/* #endif */
/*         AdjustToMinimumContentsLengthWithIcon = AdjustToContentsOnFirstShow + 2 */
/*     }; */
/*     Q_ENUM(SizeAdjustPolicy) */

/*     SizeAdjustPolicy sizeAdjustPolicy() const; */
/*     void setSizeAdjustPolicy(SizeAdjustPolicy policy); */
/*     int minimumContentsLength() const; */
/*     void setMinimumContentsLength(int characters); */
/*     QSize iconSize() const; */
/*     void setIconSize(const QSize &size); */

/*     void setPlaceholderText(const QString &placeholderText); */
/*     QString placeholderText() const; */

/*     bool isEditable() const; */
/*     void setEditable(bool editable); */
/*     void setLineEdit(QLineEdit *edit); */
/*     QLineEdit *lineEdit() const; */
/* #ifndef QT_NO_VALIDATOR */
/*     void setValidator(const QValidator *v); */
/*     const QValidator *validator() const; */
/* #endif */

/* #if QT_CONFIG(completer) */
/*     void setCompleter(QCompleter *c); */
/*     QCompleter *completer() const; */
/* #endif */

/*     QAbstractItemDelegate *itemDelegate() const; */
/*     void setItemDelegate(QAbstractItemDelegate *delegate); */

/*     QAbstractItemModel *model() const; */
/*     void setModel(QAbstractItemModel *model); */

/*     QModelIndex rootModelIndex() const; */
/*     void setRootModelIndex(const QModelIndex &index); */

/*     int modelColumn() const; */
/*     void setModelColumn(int visibleColumn); */

/*     int currentIndex() const; */
/*     QString currentText() const; */
/*     QVariant currentData(int role = Qt::UserRole) const; */

/*     QString itemText(int index) const; */
/*     QIcon itemIcon(int index) const; */
/*     QVariant itemData(int index, int role = Qt::UserRole) const; */

/*     inline void addItem(const QString &text, const QVariant &userData = QVariant()); */
/*     inline void addItem(const QIcon &icon, const QString &text, */
/*                         const QVariant &userData = QVariant()); */
/*     inline void addItems(const QStringList &texts) */
/*         { insertItems(count(), texts); } */

/*     inline void insertItem(int index, const QString &text, const QVariant &userData = QVariant()); */
/*     void insertItem(int index, const QIcon &icon, const QString &text, */
/*                     const QVariant &userData = QVariant()); */
/*     void insertItems(int index, const QStringList &texts); */
/*     void insertSeparator(int index); */

/*     void removeItem(int index); */

/*     void setItemText(int index, const QString &text); */
/*     void setItemIcon(int index, const QIcon &icon); */
/*     void setItemData(int index, const QVariant &value, int role = Qt::UserRole); */

/*     QAbstractItemView *view() const; */
/*     void setView(QAbstractItemView *itemView); */

/*     QSize sizeHint() const override; */
/*     QSize minimumSizeHint() const override; */

/*     virtual void showPopup(); */
/*     virtual void hidePopup(); */

/*     bool event(QEvent *event) override; */
/*     QVariant inputMethodQuery(Qt::InputMethodQuery) const override; */
/*     Q_INVOKABLE QVariant inputMethodQuery(Qt::InputMethodQuery query, const QVariant &argument) const; */

/* public Q_SLOTS: */
/*     void clear(); */
/*     void clearEditText(); */
/*     void setEditText(const QString &text); */
/*     void setCurrentIndex(int index); */
/*     void setCurrentText(const QString &text); */

/* Q_SIGNALS: */
/*     void editTextChanged(const QString &); */
/*     void activated(int index); */
/*     void textActivated(const QString &); */
/*     void highlighted(int index); */
/*     void textHighlighted(const QString &); */
/*     void currentIndexChanged(int index); */
/* #if QT_DEPRECATED_SINCE(5, 15) */
/*     QT_DEPRECATED_VERSION_X_5_15( */
/*             "Use currentIndexChanged(int) instead, and get the text using itemText(index)") */
/*     void currentIndexChanged(const QString &); */
/* #endif */
/*     void currentTextChanged(const QString &); */
/* #if QT_DEPRECATED_SINCE(5, 15) */
/*     QT_DEPRECATED_VERSION_X(5, 15, "Use textActivated() instead") */
/*     void activated(const QString &); */
/*     QT_DEPRECATED_VERSION_X(5, 15, "Use textHighlighted() instead") */
/*     void highlighted(const QString &); */
/* #endif */

/* protected: */
/*     void focusInEvent(QFocusEvent *e) override; */
/*     void focusOutEvent(QFocusEvent *e) override; */
/*     void changeEvent(QEvent *e) override; */
/*     void resizeEvent(QResizeEvent *e) override; */
/*     void paintEvent(QPaintEvent *e) override; */
/*     void showEvent(QShowEvent *e) override; */
/*     void hideEvent(QHideEvent *e) override; */
/*     void mousePressEvent(QMouseEvent *e) override; */
/*     void mouseReleaseEvent(QMouseEvent *e) override; */
/*     void keyPressEvent(QKeyEvent *e) override; */
/*     void keyReleaseEvent(QKeyEvent *e) override; */
/* #if QT_CONFIG(wheelevent) */
/*     void wheelEvent(QWheelEvent *e) override; */
/* #endif */
/* #ifndef QT_NO_CONTEXTMENU */
/*     void contextMenuEvent(QContextMenuEvent *e) override; */
/* #endif // QT_NO_CONTEXTMENU */
/*     void inputMethodEvent(QInputMethodEvent *) override; */
/*     void initStyleOption(QStyleOptionComboBox *option) const; */

}; // QComboBox


class xQStackedWidget : public QStackedWidget
{
 public:
    explicit xQStackedWidget(QWidget *parent = nullptr);
    //~QStackedWidget();

    int addWidget(QWidget *w);
    int insertWidget(int index, QWidget *w);
    void removeWidget(QWidget *w);

    QWidget *currentWidget() const;
    int currentIndex() const;

    int indexOf(QWidget *) const;
    QWidget *widget(int) const;
    int count() const;

    public Q_SLOTS:
        void setCurrentIndex(int index);
        void setCurrentWidget(QWidget *w);

 Q_SIGNALS:
        void currentChanged(int);
        void widgetRemoved(int index);

 /* protected: */
 /*        bool event(QEvent *e) override; */

}; // QStackedWidget


class xQGraphicsItem : public QGraphicsItem {
 public:
}; // QGraphicsItems

class xQGraphicsEffect : public QGraphicsEffect {
 public:
}; // QGraphicsEffect

class xQAbstractItemView : public QAbstractItemView
{
public:
    explicit xQAbstractItemView(QWidget *parent = nullptr);
    //    ~QAbstractItemView();

/*     virtual void setModel(QAbstractItemModel *model); */
/*     QAbstractItemModel *model() const; */

/*     virtual void setSelectionModel(QItemSelectionModel *selectionModel); */
/*     QItemSelectionModel *selectionModel() const; */

/*     void setItemDelegate(QAbstractItemDelegate *delegate); */
/*     QAbstractItemDelegate *itemDelegate() const; */

/*     void setSelectionMode(QAbstractItemView::SelectionMode mode); */
/*     QAbstractItemView::SelectionMode selectionMode() const; */

/*     void setSelectionBehavior(QAbstractItemView::SelectionBehavior behavior); */
/*     QAbstractItemView::SelectionBehavior selectionBehavior() const; */

/*     QModelIndex currentIndex() const; */
/*     QModelIndex rootIndex() const; */

/*     void setEditTriggers(EditTriggers triggers); */
/*     EditTriggers editTriggers() const; */

/*     void setVerticalScrollMode(ScrollMode mode); */
/*     ScrollMode verticalScrollMode() const; */
/*     void resetVerticalScrollMode(); */

/*     void setHorizontalScrollMode(ScrollMode mode); */
/*     ScrollMode horizontalScrollMode() const; */
/*     void resetHorizontalScrollMode(); */

/*     void setAutoScroll(bool enable); */
/*     bool hasAutoScroll() const; */

/*     void setAutoScrollMargin(int margin); */
/*     int autoScrollMargin() const; */

/*     void setTabKeyNavigation(bool enable); */
/*     bool tabKeyNavigation() const; */

/* #if QT_CONFIG(draganddrop) */
/*     void setDropIndicatorShown(bool enable); */
/*     bool showDropIndicator() const; */

/*     void setDragEnabled(bool enable); */
/*     bool dragEnabled() const; */

/*     void setDragDropOverwriteMode(bool overwrite); */
/*     bool dragDropOverwriteMode() const; */

/*     enum DragDropMode { */
/*         NoDragDrop, */
/*         DragOnly, */
/*         DropOnly, */
/*         DragDrop, */
/*         InternalMove */
/*     }; */
/*     Q_ENUM(DragDropMode) */

/*     void setDragDropMode(DragDropMode behavior); */
/*     DragDropMode dragDropMode() const; */

/*     void setDefaultDropAction(Qt::DropAction dropAction); */
/*     Qt::DropAction defaultDropAction() const; */
/* #endif */

/*     void setAlternatingRowColors(bool enable); */
/*     bool alternatingRowColors() const; */

/*     void setIconSize(const QSize &size); */
/*     QSize iconSize() const; */

/*     void setTextElideMode(Qt::TextElideMode mode); */
/*     Qt::TextElideMode textElideMode() const; */

/*     virtual void keyboardSearch(const QString &search); */

/*     virtual QRect visualRect(const QModelIndex &index) const = 0; */
/*     virtual void scrollTo(const QModelIndex &index, ScrollHint hint = EnsureVisible) = 0; */
/*     virtual QModelIndex indexAt(const QPoint &point) const = 0; */

/*     QSize sizeHintForIndex(const QModelIndex &index) const; */
/*     virtual int sizeHintForRow(int row) const; */
/*     virtual int sizeHintForColumn(int column) const; */

/*     void openPersistentEditor(const QModelIndex &index); */
/*     void closePersistentEditor(const QModelIndex &index); */
/*     bool isPersistentEditorOpen(const QModelIndex &index) const; */

/*     void setIndexWidget(const QModelIndex &index, QWidget *widget); */
/*     QWidget *indexWidget(const QModelIndex &index) const; */

/*     void setItemDelegateForRow(int row, QAbstractItemDelegate *delegate); */
/*     QAbstractItemDelegate *itemDelegateForRow(int row) const; */

/*     void setItemDelegateForColumn(int column, QAbstractItemDelegate *delegate); */
/*     QAbstractItemDelegate *itemDelegateForColumn(int column) const; */

/*     QAbstractItemDelegate *itemDelegate(const QModelIndex &index) const; */

/*     virtual QVariant inputMethodQuery(Qt::InputMethodQuery query) const override; */

/*     using QAbstractScrollArea::update; */

/* public Q_SLOTS: */
/*     virtual void reset(); */
/*     virtual void setRootIndex(const QModelIndex &index); */
/*     virtual void doItemsLayout(); */
/*     virtual void selectAll(); */
/*     void edit(const QModelIndex &index); */
/*     void clearSelection(); */
/*     void setCurrentIndex(const QModelIndex &index); */
/*     void scrollToTop(); */
/*     void scrollToBottom(); */
/*     void update(const QModelIndex &index); */

/* protected Q_SLOTS: */
/*     virtual void dataChanged(const QModelIndex &topLeft, const QModelIndex &bottomRight, const QVector<int> &roles = QVector<int>()); */
/*     virtual void rowsInserted(const QModelIndex &parent, int start, int end); */
/*     virtual void rowsAboutToBeRemoved(const QModelIndex &parent, int start, int end); */
/*     virtual void selectionChanged(const QItemSelection &selected, const QItemSelection &deselected); */
/*     virtual void currentChanged(const QModelIndex &current, const QModelIndex &previous); */
/*     virtual void updateEditorData(); */
/*     virtual void updateEditorGeometries(); */
/*     virtual void updateGeometries(); */
/*     virtual void verticalScrollbarAction(int action); */
/*     virtual void horizontalScrollbarAction(int action); */
/*     virtual void verticalScrollbarValueChanged(int value); */
/*     virtual void horizontalScrollbarValueChanged(int value); */
/*     virtual void closeEditor(QWidget *editor, QAbstractItemDelegate::EndEditHint hint); */
/*     virtual void commitData(QWidget *editor); */
/*     virtual void editorDestroyed(QObject *editor); */

/* Q_SIGNALS: */
/*     void pressed(const QModelIndex &index); */
/*     void clicked(const QModelIndex &index); */
/*     void doubleClicked(const QModelIndex &index); */

/*     void activated(const QModelIndex &index); */
/*     void entered(const QModelIndex &index); */
/*     void viewportEntered(); */

/*     void iconSizeChanged(const QSize &size); */

/* protected: */
/*     QAbstractItemView(QAbstractItemViewPrivate &, QWidget *parent = nullptr); */

/* #if QT_DEPRECATED_SINCE(5, 13) */
/*     QT_DEPRECATED void setHorizontalStepsPerItem(int steps); */
/*     QT_DEPRECATED int horizontalStepsPerItem() const; */
/*     QT_DEPRECATED void setVerticalStepsPerItem(int steps); */
/*     QT_DEPRECATED int verticalStepsPerItem() const; */
/* #endif */

/*     enum CursorAction { MoveUp, MoveDown, MoveLeft, MoveRight, */
/*                         MoveHome, MoveEnd, MovePageUp, MovePageDown, */
/*                         MoveNext, MovePrevious }; */
/*     virtual QModelIndex moveCursor(CursorAction cursorAction, */
/*                                    Qt::KeyboardModifiers modifiers) = 0; */

/*     virtual int horizontalOffset() const = 0; */
/*     virtual int verticalOffset() const = 0; */

/*     virtual bool isIndexHidden(const QModelIndex &index) const = 0; */

/*     virtual void setSelection(const QRect &rect, QItemSelectionModel::SelectionFlags command) = 0; */
/*     virtual QRegion visualRegionForSelection(const QItemSelection &selection) const = 0; */
/*     virtual QModelIndexList selectedIndexes() const; */

/*     virtual bool edit(const QModelIndex &index, EditTrigger trigger, QEvent *event); */

/*     virtual QItemSelectionModel::SelectionFlags selectionCommand(const QModelIndex &index, */
/*                                                                  const QEvent *event = nullptr) const; */

/* #if QT_CONFIG(draganddrop) */
/*     virtual void startDrag(Qt::DropActions supportedActions); */
/* #endif */

/*     virtual QStyleOptionViewItem viewOptions() const; */

/*     enum State { */
/*         NoState, */
/*         DraggingState, */
/*         DragSelectingState, */
/*         EditingState, */
/*         ExpandingState, */
/*         CollapsingState, */
/*         AnimatingState */
/*     }; */

/*     State state() const; */
/*     void setState(State state); */

/*     void scheduleDelayedItemsLayout(); */
/*     void executeDelayedItemsLayout(); */

/*     void setDirtyRegion(const QRegion &region); */
/*     void scrollDirtyRegion(int dx, int dy); */
/*     QPoint dirtyRegionOffset() const; */

/*     void startAutoScroll(); */
/*     void stopAutoScroll(); */
/*     void doAutoScroll(); */

/*     bool focusNextPrevChild(bool next) override; */
/*     bool event(QEvent *event) override; */
/*     bool viewportEvent(QEvent *event) override; */
/*     void mousePressEvent(QMouseEvent *event) override; */
/*     void mouseMoveEvent(QMouseEvent *event) override; */
/*     void mouseReleaseEvent(QMouseEvent *event) override; */
/*     void mouseDoubleClickEvent(QMouseEvent *event) override; */
/* #if QT_CONFIG(draganddrop) */
/*     void dragEnterEvent(QDragEnterEvent *event) override; */
/*     void dragMoveEvent(QDragMoveEvent *event) override; */
/*     void dragLeaveEvent(QDragLeaveEvent *event) override; */
/*     void dropEvent(QDropEvent *event) override; */
/* #endif */
/*     void focusInEvent(QFocusEvent *event) override; */
/*     void focusOutEvent(QFocusEvent *event) override; */
/*     void keyPressEvent(QKeyEvent *event) override; */
/*     void resizeEvent(QResizeEvent *event) override; */
/*     void timerEvent(QTimerEvent *event) override; */
/*     void inputMethodEvent(QInputMethodEvent *event) override; */
/*     bool eventFilter(QObject *object, QEvent *event) override; */

/* #if QT_CONFIG(draganddrop) */
/*     enum DropIndicatorPosition { OnItem, AboveItem, BelowItem, OnViewport }; */
/*     DropIndicatorPosition dropIndicatorPosition() const; */
/* #endif */

/*     QSize viewportSizeHint() const override; */

}; // QAbstractItemView

class  xQListView : public QListView
{
public:
    explicit xQListView(QWidget *parent = nullptr);
    //    ~QListView();

/*     void setMovement(Movement movement); */
/*     Movement movement() const; */

/*     void setFlow(Flow flow); */
/*     Flow flow() const; */

/*     void setWrapping(bool enable); */
/*     bool isWrapping() const; */

/*     void setResizeMode(ResizeMode mode); */
/*     ResizeMode resizeMode() const; */

/*     void setLayoutMode(LayoutMode mode); */
/*     LayoutMode layoutMode() const; */

/*     void setSpacing(int space); */
/*     int spacing() const; */

/*     void setBatchSize(int batchSize); */
/*     int batchSize() const; */

/*     void setGridSize(const QSize &size); */
/*     QSize gridSize() const; */

/*     void setViewMode(ViewMode mode); */
/*     ViewMode viewMode() const; */

/*     void clearPropertyFlags(); */

/*     bool isRowHidden(int row) const; */
/*     void setRowHidden(int row, bool hide); */

/*     void setModelColumn(int column); */
/*     int modelColumn() const; */

/*     void setUniformItemSizes(bool enable); */
/*     bool uniformItemSizes() const; */

/*     void setWordWrap(bool on); */
/*     bool wordWrap() const; */

/*     void setSelectionRectVisible(bool show); */
/*     bool isSelectionRectVisible() const; */

/*     void setItemAlignment(Qt::Alignment alignment); */
/*     Qt::Alignment itemAlignment() const; */

/*     QRect visualRect(const QModelIndex &index) const override; */
/*     void scrollTo(const QModelIndex &index, ScrollHint hint = EnsureVisible) override; */
/*     QModelIndex indexAt(const QPoint &p) const override; */

/*     void doItemsLayout() override; */
/*     void reset() override; */
/*     void setRootIndex(const QModelIndex &index) override; */

/* Q_SIGNALS: */
/*     void indexesMoved(const QModelIndexList &indexes); */

/* protected: */
/*     QListView(QListViewPrivate &, QWidget *parent = nullptr); */

/*     bool event(QEvent *e) override; */

/*     void scrollContentsBy(int dx, int dy) override; */

/*     void resizeContents(int width, int height); */
/*     QSize contentsSize() const; */

/*     void dataChanged(const QModelIndex &topLeft, const QModelIndex &bottomRight, const QVector<int> &roles = QVector<int>()) override; */
/*     void rowsInserted(const QModelIndex &parent, int start, int end) override; */
/*     void rowsAboutToBeRemoved(const QModelIndex &parent, int start, int end) override; */

/*     void mouseMoveEvent(QMouseEvent *e) override; */
/*     void mouseReleaseEvent(QMouseEvent *e) override; */
/* #if QT_CONFIG(wheelevent) */
/*     void wheelEvent(QWheelEvent *e) override; */
/* #endif */

/*     void timerEvent(QTimerEvent *e) override; */
/*     void resizeEvent(QResizeEvent *e) override; */
/* #if QT_CONFIG(draganddrop) */
/*     void dragMoveEvent(QDragMoveEvent *e) override; */
/*     void dragLeaveEvent(QDragLeaveEvent *e) override; */
/*     void dropEvent(QDropEvent *e) override; */
/*     void startDrag(Qt::DropActions supportedActions) override; */
/* #endif // QT_CONFIG(draganddrop) */

/*     QStyleOptionViewItem viewOptions() const override; */
/*     void paintEvent(QPaintEvent *e) override; */

/*     int horizontalOffset() const override; */
/*     int verticalOffset() const override; */
/*     QModelIndex moveCursor(CursorAction cursorAction, Qt::KeyboardModifiers modifiers) override; */
/*     QRect rectForIndex(const QModelIndex &index) const; */
/*     void setPositionForIndex(const QPoint &position, const QModelIndex &index); */

/*     void setSelection(const QRect &rect, QItemSelectionModel::SelectionFlags command) override; */
/*     QRegion visualRegionForSelection(const QItemSelection &selection) const override; */
/*     QModelIndexList selectedIndexes() const override; */

/*     void updateGeometries() override; */

/*     bool isIndexHidden(const QModelIndex &index) const override; */

/*     void selectionChanged(const QItemSelection &selected, const QItemSelection &deselected) override; */
/*     void currentChanged(const QModelIndex &current, const QModelIndex &previous) override; */

/*     QSize viewportSizeHint() const override; */

}; // QListView

class xQTreeView : public QTreeView
{
public:
    explicit xQTreeView(QWidget *parent = nullptr);
    //    ~QTreeView();

/*     void setModel(QAbstractItemModel *model) override; */
/*     void setRootIndex(const QModelIndex &index) override; */
/*     void setSelectionModel(QItemSelectionModel *selectionModel) override; */

/*     QHeaderView *header() const; */
/*     void setHeader(QHeaderView *header); */

/*     int autoExpandDelay() const; */
/*     void setAutoExpandDelay(int delay); */

/*     int indentation() const; */
/*     void setIndentation(int i); */
/*     void resetIndentation(); */

/*     bool rootIsDecorated() const; */
/*     void setRootIsDecorated(bool show); */

/*     bool uniformRowHeights() const; */
/*     void setUniformRowHeights(bool uniform); */

/*     bool itemsExpandable() const; */
/*     void setItemsExpandable(bool enable); */

/*     bool expandsOnDoubleClick() const; */
/*     void setExpandsOnDoubleClick(bool enable); */

/*     int columnViewportPosition(int column) const; */
/*     int columnWidth(int column) const; */
/*     void setColumnWidth(int column, int width); */
/*     int columnAt(int x) const; */

/*     bool isColumnHidden(int column) const; */
/*     void setColumnHidden(int column, bool hide); */

/*     bool isHeaderHidden() const; */
/*     void setHeaderHidden(bool hide); */

/*     bool isRowHidden(int row, const QModelIndex &parent) const; */
/*     void setRowHidden(int row, const QModelIndex &parent, bool hide); */

/*     bool isFirstColumnSpanned(int row, const QModelIndex &parent) const; */
/*     void setFirstColumnSpanned(int row, const QModelIndex &parent, bool span); */

/*     bool isExpanded(const QModelIndex &index) const; */
/*     void setExpanded(const QModelIndex &index, bool expand); */

/*     void setSortingEnabled(bool enable); */
/*     bool isSortingEnabled() const; */

/*     void setAnimated(bool enable); */
/*     bool isAnimated() const; */

/*     void setAllColumnsShowFocus(bool enable); */
/*     bool allColumnsShowFocus() const; */

/*     void setWordWrap(bool on); */
/*     bool wordWrap() const; */

/*     void setTreePosition(int logicalIndex); */
/*     int treePosition() const; */

/*     void keyboardSearch(const QString &search) override; */

/*     QRect visualRect(const QModelIndex &index) const override; */
/*     void scrollTo(const QModelIndex &index, ScrollHint hint = EnsureVisible) override; */
/*     QModelIndex indexAt(const QPoint &p) const override; */
/*     QModelIndex indexAbove(const QModelIndex &index) const; */
/*     QModelIndex indexBelow(const QModelIndex &index) const; */

/*     void doItemsLayout() override; */
/*     void reset() override; */


/*     void dataChanged(const QModelIndex &topLeft, const QModelIndex &bottomRight, const QVector<int> &roles = QVector<int>()) override; */
/*     void selectAll() override; */

/* Q_SIGNALS: */
/*     void expanded(const QModelIndex &index); */
/*     void collapsed(const QModelIndex &index); */

/* public Q_SLOTS: */
/*     void hideColumn(int column); */
/*     void showColumn(int column); */
/*     void expand(const QModelIndex &index); */
/*     void collapse(const QModelIndex &index); */
/*     void resizeColumnToContents(int column); */
/* #if QT_DEPRECATED_SINCE(5, 13) */
/*     QT_DEPRECATED_X ("Use QTreeView::sortByColumn(int column, Qt::SortOrder order) instead") */
/*     void sortByColumn(int column); */
/* #endif */
/*     void sortByColumn(int column, Qt::SortOrder order); */
/*     void expandAll(); */
/*     void expandRecursively(const QModelIndex &index, int depth = -1); */
/*     void collapseAll(); */
/*     void expandToDepth(int depth); */

/* protected Q_SLOTS: */
/*     void columnResized(int column, int oldSize, int newSize); */
/*     void columnCountChanged(int oldCount, int newCount); */
/*     void columnMoved(); */
/*     void reexpand(); */
/*     void rowsRemoved(const QModelIndex &parent, int first, int last); */
/*     void verticalScrollbarValueChanged(int value) override; */

/* protected: */
/*     QTreeView(QTreeViewPrivate &dd, QWidget *parent = nullptr); */
/*     void scrollContentsBy(int dx, int dy) override; */
/*     void rowsInserted(const QModelIndex &parent, int start, int end) override; */
/*     void rowsAboutToBeRemoved(const QModelIndex &parent, int start, int end) override; */

/*     QModelIndex moveCursor(CursorAction cursorAction, Qt::KeyboardModifiers modifiers) override; */
/*     int horizontalOffset() const override; */
/*     int verticalOffset() const override; */

/*     void setSelection(const QRect &rect, QItemSelectionModel::SelectionFlags command) override; */
/*     QRegion visualRegionForSelection(const QItemSelection &selection) const override; */
/*     QModelIndexList selectedIndexes() const override; */

/*     void timerEvent(QTimerEvent *event) override; */
/*     void paintEvent(QPaintEvent *event) override; */

/*     void drawTree(QPainter *painter, const QRegion &region) const; */
/*     virtual void drawRow(QPainter *painter, */
/*                          const QStyleOptionViewItem &options, */
/*                          const QModelIndex &index) const; */
/*     virtual void drawBranches(QPainter *painter, */
/*                               const QRect &rect, */
/*                               const QModelIndex &index) const; */

/*     void mousePressEvent(QMouseEvent *event) override; */
/*     void mouseReleaseEvent(QMouseEvent *event) override; */
/*     void mouseDoubleClickEvent(QMouseEvent *event) override; */
/*     void mouseMoveEvent(QMouseEvent *event) override; */
/*     void keyPressEvent(QKeyEvent *event) override; */
/* #if QT_CONFIG(draganddrop) */
/*     void dragMoveEvent(QDragMoveEvent *event) override; */
/* #endif */
/*     bool viewportEvent(QEvent *event) override; */

/*     void updateGeometries() override; */

/*     QSize viewportSizeHint() const override; */

/*     int sizeHintForColumn(int column) const override; */
/*     int indexRowSizeHint(const QModelIndex &index) const; */
/*     int rowHeight(const QModelIndex &index) const; */

/*     void horizontalScrollbarAction(int action) override; */

/*     bool isIndexHidden(const QModelIndex &index) const override; */
/*     void selectionChanged(const QItemSelection &selected, */
/*                           const QItemSelection &deselected) override; */
/*     void currentChanged(const QModelIndex &current, const QModelIndex &previous) override; */

}; // QTreeView


class xQSystemTrayIcon : public QSystemTrayIcon
{
public:
    xQSystemTrayIcon(QObject *parent = nullptr);
    //xQSystemTrayIcon(const QIcon &icon, QObject *parent = nullptr);
    //~QSystemTrayIcon();

/*     enum ActivationReason { */
/*         Unknown, */
/*         Context, */
/*         DoubleClick, */
/*         Trigger, */
/*         MiddleClick */
/*     }; */

/* #if QT_CONFIG(menu) */
/*     void setContextMenu(QMenu *menu); */
/*     QMenu *contextMenu() const; */
/* #endif */

/*     QIcon icon() const; */
/*     void setIcon(const QIcon &icon); */

    QString toolTip() const;
    void setToolTip(const QString &tip);

    static bool isSystemTrayAvailable();
    static bool supportsMessages();

/*     enum MessageIcon { NoIcon, Information, Warning, Critical }; */

/*     QRect geometry() const; */
    bool isVisible() const;

public Q_SLOTS:
    void setVisible(bool visible);
    inline void show() { setVisible(true); }
    inline void hide() { setVisible(false); }
/*     void showMessage(const QString &title, const QString &msg, const QIcon &icon, int msecs = 10000); */
/*     void showMessage(const QString &title, const QString &msg, */
/*                      QSystemTrayIcon::MessageIcon icon = QSystemTrayIcon::Information, int msecs = 10000); */

/* Q_SIGNALS: */
/*     void activated(QSystemTrayIcon::ActivationReason reason); */
/*     void messageClicked(); */

/* protected: */
/*     bool event(QEvent *event) override; */

}; // QSystemTrayIcon

class xQMainWindow : public QMainWindow
{
 public:
    explicit xQMainWindow(QWidget *parent = nullptr, Qt::WindowFlags flags = Qt::WindowFlags());
    // ~QMainWindow();

    QSize iconSize() const;
    void setIconSize(const QSize &iconSize);

/*     Qt::ToolButtonStyle toolButtonStyle() const; */
/*     void setToolButtonStyle(Qt::ToolButtonStyle toolButtonStyle); */

/* #if QT_CONFIG(dockwidget) */
/*     bool isAnimated() const; */
/*     bool isDockNestingEnabled() const; */
/* #endif */

/* #if QT_CONFIG(tabbar) */
/*     bool documentMode() const; */
/*     void setDocumentMode(bool enabled); */
/* #endif */

/* #if QT_CONFIG(tabwidget) */
/*     QTabWidget::TabShape tabShape() const; */
/*     void setTabShape(QTabWidget::TabShape tabShape); */
/*     QTabWidget::TabPosition tabPosition(Qt::DockWidgetArea area) const; */
/*     void setTabPosition(Qt::DockWidgetAreas areas, QTabWidget::TabPosition tabPosition); */
/* #endif // QT_CONFIG(tabwidget) */

/*     void setDockOptions(DockOptions options); */
/*     DockOptions dockOptions() const; */

/*     bool isSeparator(const QPoint &pos) const; */

/* #if QT_CONFIG(menubar) */
/*     QMenuBar *menuBar() const; */
/*     void setMenuBar(QMenuBar *menubar); */

/*     QWidget  *menuWidget() const; */
/*     void setMenuWidget(QWidget *menubar); */
/* #endif */

/* #if QT_CONFIG(statusbar) */
/*     QStatusBar *statusBar() const; */
/*     void setStatusBar(QStatusBar *statusbar); */
/* #endif */

/*     QWidget *centralWidget() const; */
/*     void setCentralWidget(QWidget *widget); */

/*     QWidget *takeCentralWidget(); */

/* #if QT_CONFIG(dockwidget) */
/*     void setCorner(Qt::Corner corner, Qt::DockWidgetArea area); */
/*     Qt::DockWidgetArea corner(Qt::Corner corner) const; */
/* #endif */

/* #if QT_CONFIG(toolbar) */
/*     void addToolBarBreak(Qt::ToolBarArea area = Qt::TopToolBarArea); */
/*     void insertToolBarBreak(QToolBar *before); */

/*     void addToolBar(Qt::ToolBarArea area, QToolBar *toolbar); */
/*     void addToolBar(QToolBar *toolbar); */
/*     QToolBar *addToolBar(const QString &title); */
/*     void insertToolBar(QToolBar *before, QToolBar *toolbar); */
/*     void removeToolBar(QToolBar *toolbar); */
/*     void removeToolBarBreak(QToolBar *before); */

/*     bool unifiedTitleAndToolBarOnMac() const; */

/*     Qt::ToolBarArea toolBarArea( */
/* #if QT_VERSION >= QT_VERSION_CHECK(6,0,0) */
/*         const */
/* #endif */
/*         QToolBar *toolbar) const; */
/*     bool toolBarBreak(QToolBar *toolbar) const; */
/* #endif */
/* #if QT_CONFIG(dockwidget) */
/*     void addDockWidget(Qt::DockWidgetArea area, QDockWidget *dockwidget); */
/*     void addDockWidget(Qt::DockWidgetArea area, QDockWidget *dockwidget, */
/*                        Qt::Orientation orientation); */
/*     void splitDockWidget(QDockWidget *after, QDockWidget *dockwidget, */
/*                          Qt::Orientation orientation); */
/* #if QT_CONFIG(tabbar) */
/*     void tabifyDockWidget(QDockWidget *first, QDockWidget *second); */
/*     QList<QDockWidget*> tabifiedDockWidgets(QDockWidget *dockwidget) const; */
/* #endif */
/*     void removeDockWidget(QDockWidget *dockwidget); */
/*     bool restoreDockWidget(QDockWidget *dockwidget); */

/*     Qt::DockWidgetArea dockWidgetArea(QDockWidget *dockwidget) const; */

/*     void resizeDocks(const QList<QDockWidget *> &docks, */
/*                      const QList<int> &sizes, Qt::Orientation orientation); */
/* #endif // QT_CONFIG(dockwidget) */

/*     QByteArray saveState(int version = 0) const; */
/*     bool restoreState(const QByteArray &state, int version = 0); */

/* #if QT_CONFIG(menu) */
/*     virtual QMenu *createPopupMenu(); */
/* #endif */

/* public Q_SLOTS: */
/* #if QT_CONFIG(dockwidget) */
/*     void setAnimated(bool enabled); */
/*     void setDockNestingEnabled(bool enabled); */
/* #endif */
/* #if QT_CONFIG(toolbar) */
/*     void setUnifiedTitleAndToolBarOnMac(bool set); */
/* #endif */

/* Q_SIGNALS: */
/*     void iconSizeChanged(const QSize &iconSize); */
/*     void toolButtonStyleChanged(Qt::ToolButtonStyle toolButtonStyle); */
/* #if QT_CONFIG(dockwidget) */
/*     void tabifiedDockWidgetActivated(QDockWidget *dockWidget); */
/* #endif */

/* protected: */
/* #ifndef QT_NO_CONTEXTMENU */
/*     void contextMenuEvent(QContextMenuEvent *event) override; */
/* #endif */
/*     bool event(QEvent *event) override; */

}; // QMainWindow

class xQApplication : public QApplication {
 public:
    xQApplication(int &argc, char **argv, int = ApplicationFlags);

/*     static QStyle *style(); */
/*     static void setStyle(QStyle*); */
/*     static QStyle *setStyle(const QString&); */
/*     enum ColorSpec { NormalColor=0, CustomColor=1, ManyColor=2 }; */
/* #if QT_DEPRECATED_SINCE(5, 8) */
/*     QT_DEPRECATED static int colorSpec(); */
/*     QT_DEPRECATED static void setColorSpec(int); */
/* #endif // QT_DEPRECATED_SINCE(5, 8) */
/* #if QT_DEPRECATED_SINCE(5, 0) */
/*     QT_DEPRECATED static inline void setGraphicsSystem(const QString &) {} */
/* #endif */

/*     using QGuiApplication::palette; */
/*     static QPalette palette(const QWidget *); */
/*     static QPalette palette(const char *className); */
/*     static void setPalette(const QPalette &, const char* className = nullptr); */
/*     static QFont font(); */
/*     static QFont font(const QWidget*); */
/*     static QFont font(const char *className); */
/*     static void setFont(const QFont &, const char* className = nullptr); */
/*     static QFontMetrics fontMetrics(); */

/* #if QT_VERSION < 0x060000 // remove these forwarders in Qt 6 */
/*     static void setWindowIcon(const QIcon &icon); */
/*     static QIcon windowIcon(); */
/* #endif */

/*     static QWidgetList allWidgets(); */
/*     static QWidgetList topLevelWidgets(); */

/*     static QDesktopWidget *desktop(); */

/*     static QWidget *activePopupWidget(); */
/*     static QWidget *activeModalWidget(); */
/*     static QWidget *focusWidget(); */

/*     static QWidget *activeWindow(); */
/*     static void setActiveWindow(QWidget* act); */

/*     static QWidget *widgetAt(const QPoint &p); */
/*     static inline QWidget *widgetAt(int x, int y) { return widgetAt(QPoint(x, y)); } */
/*     static QWidget *topLevelAt(const QPoint &p); */
/*     static inline QWidget *topLevelAt(int x, int y)  { return topLevelAt(QPoint(x, y)); } */

/*     static void beep(); */
/*     static void alert(QWidget *widget, int duration = 0); */

/*     static void setCursorFlashTime(int); */
/*     static int cursorFlashTime(); */

/*     static void setDoubleClickInterval(int); */
/*     static int doubleClickInterval(); */

/*     static void setKeyboardInputInterval(int); */
/*     static int keyboardInputInterval(); */

/* #if QT_CONFIG(wheelevent) */
/*     static void setWheelScrollLines(int); */
/*     static int wheelScrollLines(); */
/* #endif */

/*     static void setStartDragTime(int ms); */
/*     static int startDragTime(); */
/*     static void setStartDragDistance(int l); */
/*     static int startDragDistance(); */

/*     static bool isEffectEnabled(Qt::UIEffect); */
/*     static void setEffectEnabled(Qt::UIEffect, bool enable = true); */


    static int exec();
/*     bool notify(QObject *, QEvent *) override; */

/* #ifdef QT_KEYPAD_NAVIGATION */
/*     static void setNavigationMode(Qt::NavigationMode mode); */
/*     static Qt::NavigationMode navigationMode(); */
/* #endif */

/* Q_SIGNALS: */
/*     void focusChanged(QWidget *old, QWidget *now); */

/* public: */
/*     QString styleSheet() const; */
/* public Q_SLOTS: */
/* #ifndef QT_NO_STYLE_STYLESHEET */
/*     void setStyleSheet(const QString& sheet); */
/* #endif */
/*     void setAutoSipEnabled(const bool enabled); */
/*     bool autoSipEnabled() const; */
/*     static void closeAllWindows(); */
/*     static void aboutQt(); */

/* protected: */
/*     bool event(QEvent *) override; */
/*     bool compressEvent(QEvent *, QObject *receiver, QPostEventList *) override; */

}; // QApplication



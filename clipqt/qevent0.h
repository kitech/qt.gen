#include <QtGui>

class xQInputEvent : public QInputEvent
{
public:
    /* explicit QInputEvent(Type type, Qt::KeyboardModifiers modifiers = Qt::NoModifier); */
    /* ~QInputEvent(); */
    Qt::KeyboardModifiers modifiers() const;
    /* inline void setModifiers(Qt::KeyboardModifiers amodifiers) { modState = amodifiers; } */
    ulong timestamp() const;
    void setTimestamp(ulong atimestamp);
};

class  xQEnterEvent : public QEnterEvent
{
public:
/*     QEnterEvent(const QPointF &localPos, const QPointF &windowPos, const QPointF &screenPos); */
/*     ~QEnterEvent(); */

/* #ifndef QT_NO_INTEGER_EVENT_COORDINATES */
/*     inline QPoint pos() const { return l.toPoint(); } */
/*     inline QPoint globalPos() const { return s.toPoint(); } */
/*     inline int x() const { return qRound(l.x()); } */
/*     inline int y() const { return qRound(l.y()); } */
/*     inline int globalX() const { return qRound(s.x()); } */
/*     inline int globalY() const { return qRound(s.y()); } */
/* #endif */
/*     const QPointF &localPos() const { return l; } */
/*     const QPointF &windowPos() const { return w; } */
/*     const QPointF &screenPos() const { return s; } */

};

class  xQMouseEvent : public QMouseEvent
{
public:
/*     QMouseEvent(Type type, const QPointF &localPos, Qt::MouseButton button, */
/*                 Qt::MouseButtons buttons, Qt::KeyboardModifiers modifiers); */
/*     QMouseEvent(Type type, const QPointF &localPos, const QPointF &screenPos, */
/*                 Qt::MouseButton button, Qt::MouseButtons buttons, */
/*                 Qt::KeyboardModifiers modifiers); */
/*     QMouseEvent(Type type, const QPointF &localPos, const QPointF &windowPos, const QPointF &screenPos, */
/*                 Qt::MouseButton button, Qt::MouseButtons buttons, */
/*                 Qt::KeyboardModifiers modifiers); */
/*     QMouseEvent(Type type, const QPointF &localPos, const QPointF &windowPos, const QPointF &screenPos, */
/*                 Qt::MouseButton button, Qt::MouseButtons buttons, */
/*                 Qt::KeyboardModifiers modifiers, Qt::MouseEventSource source); */
/*     ~QMouseEvent(); */

#ifndef QT_NO_INTEGER_EVENT_COORDINATES
    inline QPoint pos() const;
    inline QPoint globalPos() const;
    inline int x() const;
    inline int y() const;
    inline int globalX() const;
    inline int globalY() const;
#endif
    const QPointF &localPos() const;
    const QPointF &windowPos() const;
    const QPointF &screenPos() const;

    Qt::MouseButton button() const;
    Qt::MouseButtons buttons() const;

/*     inline void setLocalPos(const QPointF &localPosition) { l = localPosition; } */

    Qt::MouseEventSource source() const;
    Qt::MouseEventFlags flags() const;

};

class  xQHoverEvent : public QHoverEvent
{
public:
/*     QHoverEvent(Type type, const QPointF &pos, const QPointF &oldPos, Qt::KeyboardModifiers modifiers = Qt::NoModifier); */
/*     ~QHoverEvent(); */

/* #ifndef QT_NO_INTEGER_EVENT_COORDINATES */
/*     inline QPoint pos() const { return p.toPoint(); } */
/*     inline QPoint oldPos() const { return op.toPoint(); } */
/* #endif */

/*     inline const QPointF &posF() const { return p; } */
/*     inline const QPointF &oldPosF() const { return op; } */
};

#if QT_CONFIG(wheelevent)
class  xQWheelEvent : public QWheelEvent
{
public:
    /* QWheelEvent(QPointF pos, QPointF globalPos, QPoint pixelDelta, QPoint angleDelta, */
    /*             Qt::MouseButtons buttons, Qt::KeyboardModifiers modifiers, Qt::ScrollPhase phase, */
    /*             bool inverted, Qt::MouseEventSource source = Qt::MouseEventNotSynthesized); */
    /* ~QWheelEvent(); */


    inline QPoint pixelDelta() const;
    inline QPoint angleDelta() const;

    inline QPointF position() const;
    inline QPointF globalPosition() const;

    Qt::MouseButtons buttons() const;

    Qt::ScrollPhase phase() const;
    bool inverted() const;

    Qt::MouseEventSource source() const;
};
#endif

#if QT_CONFIG(tabletevent)
class  xQTabletEvent : public QTabletEvent
{
public:
    /* QTabletEvent(Type t, const QPointF &pos, const QPointF &globalPos, */
    /*              int device, int pointerType, qreal pressure, int xTilt, int yTilt, */
    /*              qreal tangentialPressure, qreal rotation, int z, */
    /*              Qt::KeyboardModifiers keyState, qint64 uniqueID, */
    /*              Qt::MouseButton button, Qt::MouseButtons buttons); */
    /* ~QTabletEvent(); */

    /* inline QPoint pos() const { return mPos.toPoint(); } */
    /* inline QPoint globalPos() const { return mGPos.toPoint(); } */

    /* inline const QPointF &posF() const { return mPos; } */
    /* inline const QPointF &globalPosF() const { return mGPos; } */

    inline int x() const;
    inline int y() const;
    inline int globalX() const ;
    inline int globalY() const;
    /* inline TabletDevice deviceType() const { return TabletDevice(mDev); } */
    /* inline PointerType pointerType() const { return PointerType(mPointerType); } */
    /* inline qint64 uniqueId() const { return mUnique; } */
    /* inline qreal pressure() const { return mPress; } */
    /* inline int z() const { return mZ; } */
    /* inline qreal tangentialPressure() const { return mTangential; } */
    /* inline qreal rotation() const { return mRot; } */
    /* inline int xTilt() const { return mXT; } */
    /* inline int yTilt() const { return mYT; } */
    /* Qt::MouseButton button() const; */
    /* Qt::MouseButtons buttons() const; */
};
#endif // QT_CONFIG(tabletevent)

#ifndef QT_NO_GESTURES
class  xQNativeGestureEvent : public QNativeGestureEvent
{
public:
/*     QNativeGestureEvent(Qt::NativeGestureType type, const QTouchDevice *dev, const QPointF &localPos, const QPointF &windowPos, */
/*                         const QPointF &screenPos, qreal value, ulong sequenceId, quint64 intArgument); */
/*     ~QNativeGestureEvent(); */
/*     Qt::NativeGestureType gestureType() const { return mGestureType; } */
/*     qreal value() const { return mRealValue; } */

/* #ifndef QT_NO_INTEGER_EVENT_COORDINATES */
/*     inline const QPoint pos() const { return mLocalPos.toPoint(); } */
/*     inline const QPoint globalPos() const { return mScreenPos.toPoint(); } */
/* #endif */
/*     const QPointF &localPos() const { return mLocalPos; } */
/*     const QPointF &windowPos() const { return mWindowPos; } */
/*     const QPointF &screenPos() const { return mScreenPos; } */

/*     const QTouchDevice *device() const; */
};
#endif // QT_NO_GESTURES

class  xQKeyEvent : public QKeyEvent
{
public:
/*     QKeyEvent(Type type, int key, Qt::KeyboardModifiers modifiers, const QString& text = QString(), */
/*               bool autorep = false, ushort count = 1); */
/*     QKeyEvent(Type type, int key, Qt::KeyboardModifiers modifiers, */
/*               quint32 nativeScanCode, quint32 nativeVirtualKey, quint32 nativeModifiers, */
/*               const QString &text = QString(), bool autorep = false, ushort count = 1); */
/*     ~QKeyEvent(); */

    int key() const;
/* #ifndef QT_NO_SHORTCUT */
/*     bool matches(QKeySequence::StandardKey key) const; */
/* #endif */
    Qt::KeyboardModifiers modifiers() const;
    inline QString text() const;
    inline bool isAutoRepeat() const;
    inline int count() const;

    inline quint32 nativeScanCode() const;
    inline quint32 nativeVirtualKey() const;
    inline quint32 nativeModifiers() const;

};


class  xQFocusEvent : public QFocusEvent
{
public:
    /* explicit QFocusEvent(Type type, Qt::FocusReason reason=Qt::OtherFocusReason); */
    /* ~QFocusEvent(); */

    /* inline bool gotFocus() const { return type() == FocusIn; } */
    /* inline bool lostFocus() const { return type() == FocusOut; } */

    Qt::FocusReason reason() const;
};


class  xQPaintEvent : public QPaintEvent
{
public:
    /* explicit QPaintEvent(const QRegion& paintRegion); */
    /* explicit QPaintEvent(const QRect &paintRect); */
    /* ~QPaintEvent(); */

    /* inline const QRect &rect() const { return m_rect; } */
    /* inline const QRegion &region() const { return m_region; } */

};

class  xQMoveEvent : public QMoveEvent
{
public:
    /* QMoveEvent(const QPoint &pos, const QPoint &oldPos); */
    /* ~QMoveEvent(); */

    /* inline const QPoint &pos() const { return p; } */
    /* inline const QPoint &oldPos() const { return oldp;} */
};

class  xQExposeEvent : public QExposeEvent
{
public:
    /* explicit QExposeEvent(const QRegion &rgn); */
    /* ~QExposeEvent(); */

    /* inline const QRegion &region() const { return rgn; } */

};

class  xQPlatformSurfaceEvent : public QPlatformSurfaceEvent
{
public:
    /* explicit QPlatformSurfaceEvent(SurfaceEventType surfaceEventType); */
    /* ~QPlatformSurfaceEvent(); */

    /* inline SurfaceEventType surfaceEventType() const { return m_surfaceEventType; } */
};

class  xQResizeEvent : public QResizeEvent
{
public:
    /* QResizeEvent(const QSize &size, const QSize &oldSize); */
    /* ~QResizeEvent(); */

    /* inline const QSize &size() const { return s; } */
    /* inline const QSize &oldSize()const { return olds;} */
};


class  xQCloseEvent : public QCloseEvent
{
public:
    /* QCloseEvent(); */
    /* ~QCloseEvent(); */
};


class  xQIconDragEvent : public QIconDragEvent
{
public:
    /* QIconDragEvent(); */
    /* ~QIconDragEvent(); */
};


class  xQShowEvent : public QShowEvent
{
public:
    /* QShowEvent(); */
    /* ~QShowEvent(); */
};


class  xQHideEvent : public QHideEvent
{
public:
    /* QHideEvent(); */
    /* ~QHideEvent(); */
};

#ifndef QT_NO_CONTEXTMENU
class  xQContextMenuEvent : public QContextMenuEvent
{
public:
    /* QContextMenuEvent(Reason reason, const QPoint &pos, const QPoint &globalPos, */
    /*                   Qt::KeyboardModifiers modifiers); */
    /* QContextMenuEvent(Reason reason, const QPoint &pos, const QPoint &globalPos); */
    /* QContextMenuEvent(Reason reason, const QPoint &pos); */
    /* ~QContextMenuEvent(); */

    /* inline int x() const { return p.x(); } */
    /* inline int y() const { return p.y(); } */
    /* inline int globalX() const { return gp.x(); } */
    /* inline int globalY() const { return gp.y(); } */

    /* inline const QPoint& pos() const { return p; } */
    /* inline const QPoint& globalPos() const { return gp; } */

    /* inline Reason reason() const { return Reason(reas); } */

};
#endif // QT_NO_CONTEXTMENU

#ifndef QT_NO_INPUTMETHOD
class  xQInputMethodEvent : public QInputMethodEvent
{
public:
    /* QInputMethodEvent(); */
    /* QInputMethodEvent(const QString &preeditText, const QList<Attribute> &attributes); */
    /* ~QInputMethodEvent(); */

    /* void setCommitString(const QString &commitString, int replaceFrom = 0, int replaceLength = 0); */
    /* inline const QList<Attribute> &attributes() const { return attrs; } */
    /* inline const QString &preeditString() const { return preedit; } */

    /* inline const QString &commitString() const { return commit; } */
    /* inline int replacementStart() const { return replace_from; } */
    /* inline int replacementLength() const { return replace_length; } */

    /* QInputMethodEvent(const QInputMethodEvent &other); */

};

class  xQInputMethodQueryEvent : public QInputMethodQueryEvent
{
public:
    /* explicit QInputMethodQueryEvent(Qt::InputMethodQueries queries); */
    /* ~QInputMethodQueryEvent(); */

    /* Qt::InputMethodQueries queries() const { return m_queries; } */

    /* void setValue(Qt::InputMethodQuery query, const QVariant &value); */
    /* QVariant value(Qt::InputMethodQuery query) const; */
};

#endif // QT_NO_INPUTMETHOD

#if QT_CONFIG(draganddrop)

class  xQDropEvent : public QDropEvent
{
public:
    /* QDropEvent(const QPointF& pos, Qt::DropActions actions, const QMimeData *data, */
    /*            Qt::MouseButtons buttons, Qt::KeyboardModifiers modifiers, Type type = Drop); */
    /* ~QDropEvent(); */

    /* inline QPoint pos() const { return p.toPoint(); } */
    /* inline const QPointF &posF() const { return p; } */
    /* inline Qt::MouseButtons mouseButtons() const { return mouseState; } */
    /* inline Qt::KeyboardModifiers keyboardModifiers() const { return modState; } */

    /* inline Qt::DropActions possibleActions() const { return act; } */
    /* inline Qt::DropAction proposedAction() const { return default_action; } */
    /* inline void acceptProposedAction() { drop_action = default_action; accept(); } */

    /* inline Qt::DropAction dropAction() const { return drop_action; } */
    /* void setDropAction(Qt::DropAction action); */

    /* QObject* source() const; */
    /* inline const QMimeData *mimeData() const { return mdata; } */

};


class  xQDragMoveEvent : public QDragMoveEvent
{
public:
    /* QDragMoveEvent(const QPoint &pos, Qt::DropActions actions, const QMimeData *data, */
    /*                Qt::MouseButtons buttons, Qt::KeyboardModifiers modifiers, Type type = DragMove); */
    /* ~QDragMoveEvent(); */

    /* inline QRect answerRect() const { return rect; } */

    /* inline void accept() { QDropEvent::accept(); } */
    /* inline void ignore() { QDropEvent::ignore(); } */

    /* inline void accept(const QRect & r) { accept(); rect = r; } */
    /* inline void ignore(const QRect & r) { ignore(); rect = r; } */

};


class  xQDragEnterEvent : public QDragEnterEvent
{
public:
    /* QDragEnterEvent(const QPoint &pos, Qt::DropActions actions, const QMimeData *data, */
    /*                 Qt::MouseButtons buttons, Qt::KeyboardModifiers modifiers); */
    /* ~QDragEnterEvent(); */
};


class  xQDragLeaveEvent : public QDragLeaveEvent
{
public:
    /* QDragLeaveEvent(); */
    /* ~QDragLeaveEvent(); */
};
#endif // QT_CONFIG(draganddrop)


class  xQHelpEvent : public QHelpEvent
{
public:
    /* QHelpEvent(Type type, const QPoint &pos, const QPoint &globalPos); */
    /* ~QHelpEvent(); */

    /* inline int x() const { return p.x(); } */
    /* inline int y() const { return p.y(); } */
    /* inline int globalX() const { return gp.x(); } */
    /* inline int globalY() const { return gp.y(); } */

    /* inline const QPoint& pos()  const { return p; } */
    /* inline const QPoint& globalPos() const { return gp; } */
};

#ifndef QT_NO_STATUSTIP
class  xQStatusTipEvent : public QStatusTipEvent
{
public:
    /* explicit QStatusTipEvent(const QString &tip); */
    /* ~QStatusTipEvent(); */

    /* inline QString tip() const { return s; } */
};
#endif

#if QT_CONFIG(whatsthis)
class  xQWhatsThisClickedEvent : public QWhatsThisClickedEvent
{
public:
    /* explicit QWhatsThisClickedEvent(const QString &href); */
    /* ~QWhatsThisClickedEvent(); */

    /* inline QString href() const { return s; } */
};
#endif

#ifndef QT_NO_ACTION
class  xQActionEvent : public QActionEvent
{
public:
    /* QActionEvent(int type, QAction *action, QAction *before = nullptr); */
    /* ~QActionEvent(); */

    /* inline QAction *action() const { return act; } */
    /* inline QAction *before() const { return bef; } */
};
#endif

class  xQFileOpenEvent : public QFileOpenEvent
{
public:
    /* explicit QFileOpenEvent(const QString &file); */
    /* explicit QFileOpenEvent(const QUrl &url); */
    /* ~QFileOpenEvent(); */

    inline QString file() const;
    /* QUrl url() const { return m_url; } */
    /* bool openFile(QFile &file, QIODevice::OpenMode flags) const; */
};

#ifndef QT_NO_TOOLBAR
class  xQToolBarChangeEvent : public QToolBarChangeEvent
{
public:
    /* explicit QToolBarChangeEvent(bool t); */
    /* ~QToolBarChangeEvent(); */

    /* inline bool toggle() const { return tog; } */
};
#endif

#ifndef QT_NO_SHORTCUT
class  xQShortcutEvent : public QShortcutEvent
{
public:
    /* QShortcutEvent(const QKeySequence &key, int id, bool ambiguous = false); */
    /* ~QShortcutEvent(); */

    /* inline const QKeySequence &key() const { return sequence; } */
    /* inline int shortcutId() const { return sid; } */
    /* inline bool isAmbiguous() const { return ambig; } */
};
#endif

class  xQWindowStateChangeEvent: public QWindowStateChangeEvent
{
public:
    /* explicit QWindowStateChangeEvent(Qt::WindowStates aOldState, bool isOverride = false); */
    /* ~QWindowStateChangeEvent(); */

    /* inline Qt::WindowStates oldState() const { return ostate; } */
    /* bool isOverride() const; */

};

class  xQTouchEvent : public QTouchEvent
{
public:
    /* explicit QTouchEvent(QEvent::Type eventType, */
    /*                      QTouchDevice *device = nullptr, */
    /*                      Qt::KeyboardModifiers modifiers = Qt::NoModifier, */
    /*                      Qt::TouchPointStates touchPointStates = Qt::TouchPointStates(), */
    /*                      const QList<QTouchEvent::TouchPoint> &touchPoints = QList<QTouchEvent::TouchPoint>()); */
    /* ~QTouchEvent(); */

    /* inline QWindow *window() const { return _window; } */
    /* inline QObject *target() const { return _target; } */
    /* inline Qt::TouchPointStates touchPointStates() const { return _touchPointStates; } */
    /* inline const QList<QTouchEvent::TouchPoint> &touchPoints() const { return _touchPoints; } */
    /* inline QTouchDevice *device() const { return _device; } */

    /* // internal */
    /* inline void setWindow(QWindow *awindow) { _window = awindow; } */
    /* inline void setTarget(QObject *atarget) { _target = atarget; } */
    /* inline void setTouchPointStates(Qt::TouchPointStates aTouchPointStates) { _touchPointStates = aTouchPointStates; } */
    /* inline void setTouchPoints(const QList<QTouchEvent::TouchPoint> &atouchPoints) { _touchPoints = atouchPoints; } */
    /* inline void setDevice(QTouchDevice *adevice) { _device = adevice; } */

};

class  xQScrollPrepareEvent : public QScrollPrepareEvent
{
public:
    /* explicit QScrollPrepareEvent(const QPointF &startPos); */
    /* ~QScrollPrepareEvent(); */

    /* QPointF startPos() const; */

    /* QSizeF viewportSize() const; */
    /* QRectF contentPosRange() const; */
    /* QPointF contentPos() const; */

    /* void setViewportSize(const QSizeF &size); */
    /* void setContentPosRange(const QRectF &rect); */
    /* void setContentPos(const QPointF &pos); */

};


class  xQScrollEvent : public QScrollEvent
{
public:
    /* QScrollEvent(const QPointF &contentPos, const QPointF &overshoot, ScrollState scrollState); */
    /* ~QScrollEvent(); */

    /* QPointF contentPos() const; */
    /* QPointF overshootDistance() const; */
    /* ScrollState scrollState() const; */
};

class  xQScreenOrientationChangeEvent : public QScreenOrientationChangeEvent
{
public:
    /* QScreenOrientationChangeEvent(QScreen *screen, Qt::ScreenOrientation orientation); */
    /* ~QScreenOrientationChangeEvent(); */

    /* QScreen *screen() const; */
    /* Qt::ScreenOrientation orientation() const; */
};

class  xQApplicationStateChangeEvent : public QApplicationStateChangeEvent
{
public:
    /* explicit QApplicationStateChangeEvent(Qt::ApplicationState state); */
    /* Qt::ApplicationState applicationState() const; */
};

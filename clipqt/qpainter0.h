#include <QtGui>

class xQPaintDevice : public QPaintDevice  // device for QPainter
{
public:

    virtual int devType() const;
    bool paintingActive() const;
    virtual QPaintEngine *paintEngine() const = 0;

    int width() const ;
    int height() const ;
    int widthMM() const;
    int heightMM() const;
    int logicalDpiX() const;
    int logicalDpiY() const ;
    int physicalDpiX() const ;
    int physicalDpiY() const ;
    int devicePixelRatio() const;
    qreal devicePixelRatioF()  const ;
    int colorCount() const ;
    int depth() const ;

    static inline qreal devicePixelRatioFScale();
};

class xQPaintEngine : public QPaintEngine
{
 public:

};

class xQPainter : public QPainter
{
 public:
    explicit xQPainter(QPaintDevice *);
    // ~xQPainter();

    QPaintDevice *device() const;

/*     bool begin(QPaintDevice *); */
/*     bool end(); */
/*     bool isActive() const; */

/*     void setCompositionMode(CompositionMode mode); */
/*     CompositionMode compositionMode() const; */

/*     const QFont &font() const; */
/*     void setFont(const QFont &f); */

/*     QFontMetrics fontMetrics() const; */
/*     QFontInfo fontInfo() const; */

/*     void setPen(const QColor &color); */
/*     void setPen(const QPen &pen); */
/*     void setPen(Qt::PenStyle style); */
/*     const QPen &pen() const; */

/*     void setBrush(const QBrush &brush); */
/*     void setBrush(Qt::BrushStyle style); */
/*     const QBrush &brush() const; */

/*     // attributes/modes */
/*     void setBackgroundMode(Qt::BGMode mode); */
/*     Qt::BGMode backgroundMode() const; */

/*     QPoint brushOrigin() const; */
/*     inline void setBrushOrigin(int x, int y); */
/*     inline void setBrushOrigin(const QPoint &); */
/*     void setBrushOrigin(const QPointF &); */

/*     void setBackground(const QBrush &bg); */
/*     const QBrush &background() const; */

/*     qreal opacity() const; */
/*     void setOpacity(qreal opacity); */

/*     // Clip functions */
/*     QRegion clipRegion() const; */
/*     QPainterPath clipPath() const; */

/*     void setClipRect(const QRectF &, Qt::ClipOperation op = Qt::ReplaceClip); */
/*     void setClipRect(const QRect &, Qt::ClipOperation op = Qt::ReplaceClip); */
/*     inline void setClipRect(int x, int y, int w, int h, Qt::ClipOperation op = Qt::ReplaceClip); */

/*     void setClipRegion(const QRegion &, Qt::ClipOperation op = Qt::ReplaceClip); */

/*     void setClipPath(const QPainterPath &path, Qt::ClipOperation op = Qt::ReplaceClip); */

/*     void setClipping(bool enable); */
/*     bool hasClipping() const; */

/*     QRectF clipBoundingRect() const; */

/*     void save(); */
/*     void restore(); */

/*     // XForm functions */

/*     void setTransform(const QTransform &transform, bool combine = false); */
/*     const QTransform &transform() const; */
/*     const QTransform &deviceTransform() const; */
/*     void resetTransform(); */

/*     void setWorldTransform(const QTransform &matrix, bool combine = false); */
/*     const QTransform &worldTransform() const; */

/*     QTransform combinedTransform() const; */

/*     void setWorldMatrixEnabled(bool enabled); */
/*     bool worldMatrixEnabled() const; */

/*     void scale(qreal sx, qreal sy); */
/*     void shear(qreal sh, qreal sv); */
/*     void rotate(qreal a); */

/*     void translate(const QPointF &offset); */
/*     inline void translate(const QPoint &offset); */
/*     inline void translate(qreal dx, qreal dy); */

/*     QRect window() const; */
/*     void setWindow(const QRect &window); */
/*     inline void setWindow(int x, int y, int w, int h); */

/*     QRect viewport() const; */
/*     void setViewport(const QRect &viewport); */
/*     inline void setViewport(int x, int y, int w, int h); */

/*     void setViewTransformEnabled(bool enable); */
/*     bool viewTransformEnabled() const; */

/*     // drawing functions */
/*     void strokePath(const QPainterPath &path, const QPen &pen); */
/*     void fillPath(const QPainterPath &path, const QBrush &brush); */
/*     void drawPath(const QPainterPath &path); */

/*     inline void drawPoint(const QPointF &pt); */
/*     inline void drawPoint(const QPoint &p); */
/*     inline void drawPoint(int x, int y); */

/*     void drawPoints(const QPointF *points, int pointCount); */
/*     inline void drawPoints(const QPolygonF &points); */
/*     void drawPoints(const QPoint *points, int pointCount); */
/*     inline void drawPoints(const QPolygon &points); */

/*     inline void drawLine(const QLineF &line); */
/*     inline void drawLine(const QLine &line); */
/*     inline void drawLine(int x1, int y1, int x2, int y2); */
/*     inline void drawLine(const QPoint &p1, const QPoint &p2); */
/*     inline void drawLine(const QPointF &p1, const QPointF &p2); */

/*     void drawLines(const QLineF *lines, int lineCount); */
/*     inline void drawLines(const QVector<QLineF> &lines); */
/*     void drawLines(const QPointF *pointPairs, int lineCount); */
/*     inline void drawLines(const QVector<QPointF> &pointPairs); */
/*     void drawLines(const QLine *lines, int lineCount); */
/*     inline void drawLines(const QVector<QLine> &lines); */
/*     void drawLines(const QPoint *pointPairs, int lineCount); */
/*     inline void drawLines(const QVector<QPoint> &pointPairs); */

/*     inline void drawRect(const QRectF &rect); */
/*     inline void drawRect(int x1, int y1, int w, int h); */
/*     inline void drawRect(const QRect &rect); */

/*     void drawRects(const QRectF *rects, int rectCount); */
/*     inline void drawRects(const QVector<QRectF> &rectangles); */
/*     void drawRects(const QRect *rects, int rectCount); */
/*     inline void drawRects(const QVector<QRect> &rectangles); */

/*     void drawEllipse(const QRectF &r); */
/*     void drawEllipse(const QRect &r); */
/*     inline void drawEllipse(int x, int y, int w, int h); */

/*     inline void drawEllipse(const QPointF &center, qreal rx, qreal ry); */
/*     inline void drawEllipse(const QPoint &center, int rx, int ry); */

/*     void drawPolyline(const QPointF *points, int pointCount); */
/*     inline void drawPolyline(const QPolygonF &polyline); */
/*     void drawPolyline(const QPoint *points, int pointCount); */
/*     inline void drawPolyline(const QPolygon &polygon); */

/*     void drawPolygon(const QPointF *points, int pointCount, Qt::FillRule fillRule = Qt::OddEvenFill); */
/*     inline void drawPolygon(const QPolygonF &polygon, Qt::FillRule fillRule = Qt::OddEvenFill); */
/*     void drawPolygon(const QPoint *points, int pointCount, Qt::FillRule fillRule = Qt::OddEvenFill); */
/*     inline void drawPolygon(const QPolygon &polygon, Qt::FillRule fillRule = Qt::OddEvenFill); */

/*     void drawConvexPolygon(const QPointF *points, int pointCount); */
/*     inline void drawConvexPolygon(const QPolygonF &polygon); */
/*     void drawConvexPolygon(const QPoint *points, int pointCount); */
/*     inline void drawConvexPolygon(const QPolygon &polygon); */

/*     void drawArc(const QRectF &rect, int a, int alen); */
/*     inline void drawArc(const QRect &, int a, int alen); */
/*     inline void drawArc(int x, int y, int w, int h, int a, int alen); */

/*     void drawPie(const QRectF &rect, int a, int alen); */
/*     inline void drawPie(int x, int y, int w, int h, int a, int alen); */
/*     inline void drawPie(const QRect &, int a, int alen); */

/*     void drawChord(const QRectF &rect, int a, int alen); */
/*     inline void drawChord(int x, int y, int w, int h, int a, int alen); */
/*     inline void drawChord(const QRect &, int a, int alen); */

/*     void drawRoundedRect(const QRectF &rect, qreal xRadius, qreal yRadius, */
/*                          Qt::SizeMode mode = Qt::AbsoluteSize); */
/*     inline void drawRoundedRect(int x, int y, int w, int h, qreal xRadius, qreal yRadius, */
/*                                 Qt::SizeMode mode = Qt::AbsoluteSize); */
/*     inline void drawRoundedRect(const QRect &rect, qreal xRadius, qreal yRadius, */
/*                                 Qt::SizeMode mode = Qt::AbsoluteSize); */


/*     void drawTiledPixmap(const QRectF &rect, const QPixmap &pm, const QPointF &offset = QPointF()); */
/*     inline void drawTiledPixmap(int x, int y, int w, int h, const QPixmap &, int sx=0, int sy=0); */
/*     inline void drawTiledPixmap(const QRect &, const QPixmap &, const QPoint & = QPoint()); */
/* #ifndef QT_NO_PICTURE */
/*     void drawPicture(const QPointF &p, const QPicture &picture); */
/*     inline void drawPicture(int x, int y, const QPicture &picture); */
/*     inline void drawPicture(const QPoint &p, const QPicture &picture); */
/* #endif */

/*     void drawPixmap(const QRectF &targetRect, const QPixmap &pixmap, const QRectF &sourceRect); */
/*     inline void drawPixmap(const QRect &targetRect, const QPixmap &pixmap, const QRect &sourceRect); */
/*     inline void drawPixmap(int x, int y, int w, int h, const QPixmap &pm, */
/*                            int sx, int sy, int sw, int sh); */
/*     inline void drawPixmap(int x, int y, const QPixmap &pm, */
/*                            int sx, int sy, int sw, int sh); */
/*     inline void drawPixmap(const QPointF &p, const QPixmap &pm, const QRectF &sr); */
/*     inline void drawPixmap(const QPoint &p, const QPixmap &pm, const QRect &sr); */
/*     void drawPixmap(const QPointF &p, const QPixmap &pm); */
/*     inline void drawPixmap(const QPoint &p, const QPixmap &pm); */
/*     inline void drawPixmap(int x, int y, const QPixmap &pm); */
/*     inline void drawPixmap(const QRect &r, const QPixmap &pm); */
/*     inline void drawPixmap(int x, int y, int w, int h, const QPixmap &pm); */

/*     void drawPixmapFragments(const PixmapFragment *fragments, int fragmentCount, */
/*                              const QPixmap &pixmap, PixmapFragmentHints hints = PixmapFragmentHints()); */

/*     void drawImage(const QRectF &targetRect, const QImage &image, const QRectF &sourceRect, */
/*                    Qt::ImageConversionFlags flags = Qt::AutoColor); */
/*     inline void drawImage(const QRect &targetRect, const QImage &image, const QRect &sourceRect, */
/*                           Qt::ImageConversionFlags flags = Qt::AutoColor); */
/*     inline void drawImage(const QPointF &p, const QImage &image, const QRectF &sr, */
/*                           Qt::ImageConversionFlags flags = Qt::AutoColor); */
/*     inline void drawImage(const QPoint &p, const QImage &image, const QRect &sr, */
/*                           Qt::ImageConversionFlags flags = Qt::AutoColor); */
/*     inline void drawImage(const QRectF &r, const QImage &image); */
/*     inline void drawImage(const QRect &r, const QImage &image); */
/*     void drawImage(const QPointF &p, const QImage &image); */
/*     inline void drawImage(const QPoint &p, const QImage &image); */
/*     inline void drawImage(int x, int y, const QImage &image, int sx = 0, int sy = 0, */
/*                           int sw = -1, int sh = -1, Qt::ImageConversionFlags flags = Qt::AutoColor); */

/*     void setLayoutDirection(Qt::LayoutDirection direction); */
/*     Qt::LayoutDirection layoutDirection() const; */

/* #if !defined(QT_NO_RAWFONT) */
/*     void drawGlyphRun(const QPointF &position, const QGlyphRun &glyphRun); */
/* #endif */

/*     void drawStaticText(const QPointF &topLeftPosition, const QStaticText &staticText); */
/*     inline void drawStaticText(const QPoint &topLeftPosition, const QStaticText &staticText); */
/*     inline void drawStaticText(int left, int top, const QStaticText &staticText); */

/*     void drawText(const QPointF &p, const QString &s); */
/*     inline void drawText(const QPoint &p, const QString &s); */
/*     inline void drawText(int x, int y, const QString &s); */

/*     void drawText(const QPointF &p, const QString &str, int tf, int justificationPadding); */

/*     void drawText(const QRectF &r, int flags, const QString &text, QRectF *br = nullptr); */
/*     void drawText(const QRect &r, int flags, const QString &text, QRect *br = nullptr); */
/*     inline void drawText(int x, int y, int w, int h, int flags, const QString &text, QRect *br = nullptr); */

/*     void drawText(const QRectF &r, const QString &text, const QTextOption &o = QTextOption()); */

/*     QRectF boundingRect(const QRectF &rect, int flags, const QString &text); */
/*     QRect boundingRect(const QRect &rect, int flags, const QString &text); */
/*     inline QRect boundingRect(int x, int y, int w, int h, int flags, const QString &text); */

/*     QRectF boundingRect(const QRectF &rect, const QString &text, const QTextOption &o = QTextOption()); */

/*     void drawTextItem(const QPointF &p, const QTextItem &ti); */
/*     inline void drawTextItem(int x, int y, const QTextItem &ti); */
/*     inline void drawTextItem(const QPoint &p, const QTextItem &ti); */

/*     void fillRect(const QRectF &, const QBrush &); */
/*     inline void fillRect(int x, int y, int w, int h, const QBrush &); */
/*     void fillRect(const QRect &, const QBrush &); */

/*     void fillRect(const QRectF &, const QColor &color); */
/*     inline void fillRect(int x, int y, int w, int h, const QColor &color); */
/*     void fillRect(const QRect &, const QColor &color); */

/*     inline void fillRect(int x, int y, int w, int h, Qt::GlobalColor c); */
/*     inline void fillRect(const QRect &r, Qt::GlobalColor c); */
/*     inline void fillRect(const QRectF &r, Qt::GlobalColor c); */

/*     inline void fillRect(int x, int y, int w, int h, Qt::BrushStyle style); */
/*     inline void fillRect(const QRect &r, Qt::BrushStyle style); */
/*     inline void fillRect(const QRectF &r, Qt::BrushStyle style); */

/*     inline void fillRect(int x, int y, int w, int h, QGradient::Preset preset); */
/*     inline void fillRect(const QRect &r, QGradient::Preset preset); */
/*     inline void fillRect(const QRectF &r, QGradient::Preset preset); */

/*     void eraseRect(const QRectF &); */
/*     inline void eraseRect(int x, int y, int w, int h); */
/*     inline void eraseRect(const QRect &); */

/*     void setRenderHint(RenderHint hint, bool on = true); */
/*     void setRenderHints(RenderHints hints, bool on = true); */
/*     RenderHints renderHints() const; */
/*     inline bool testRenderHint(RenderHint hint) const { return renderHints() & hint; } */

/*     QPaintEngine *paintEngine() const; */

/*     void beginNativePainting(); */
/*     void endNativePainting(); */

};

class  xQFontInfo : public QFontInfo
{
 public:

};
class  xQFontMetrics : public QFontMetrics
{
 public:

};

class  xQFont:public QFont
{
 public:
    xQFont(const QString &family, int pointSize = -1, int weight = -1, bool italic = false);
    //~xQFont();

    QString family() const;
    void setFamily(const QString &);

    //QStringList families() const;
    // void setFamilies(const QStringList &);

    QString styleName() const;
    void setStyleName(const QString &);

    int pointSize() const;
    void setPointSize(int);
    qreal pointSizeF() const;
    void setPointSizeF(qreal);

    int pixelSize() const;
    void setPixelSize(int);

    int weight() const;
    void setWeight(int);

    inline bool bold() const;
    inline void setBold(bool);

    void setStyle(Style style);
    Style style() const;

    inline bool italic() const;
    inline void setItalic(bool b);

    bool underline() const;
    void setUnderline(bool);

    bool overline() const;
    void setOverline(bool);

    bool strikeOut() const;
    void setStrikeOut(bool);

    bool fixedPitch() const;
    void setFixedPitch(bool);

    bool kerning() const;
    void setKerning(bool);

    StyleHint styleHint() const;
    StyleStrategy styleStrategy() const;
    void setStyleHint(StyleHint, StyleStrategy = PreferDefault);
    void setStyleStrategy(StyleStrategy s);

    int stretch() const;
    void setStretch(int);

    qreal letterSpacing() const;
    SpacingType letterSpacingType() const;
    void setLetterSpacing(SpacingType type, qreal spacing);

    qreal wordSpacing() const;
    void setWordSpacing(qreal spacing);

    void setCapitalization(Capitalization);
    Capitalization capitalization() const;

    void setHintingPreference(HintingPreference hintingPreference);
    HintingPreference hintingPreference() const;

    // dupicated from QFontInfo
    bool exactMatch() const;

};

class xQColor : public QColor
{
 public:

};

class xQWindow : public QWindow
{
 public:

};

class xQSurface : public QSurface {
 public:
};


class xQGuiApplication : public QGuiApplication {
 public:
};


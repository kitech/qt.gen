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

class xQPen : public QPen
{
public:
    xQPen();
    //xQPen(Qt::PenStyle);
    /* QPen(const QColor &color); */
    /* QPen(const QBrush &brush, qreal width, Qt::PenStyle s = Qt::SolidLine, */
    /*      Qt::PenCapStyle c = Qt::SquareCap, Qt::PenJoinStyle j = Qt::BevelJoin); */
    /* QPen(const QPen &pen) noexcept; */

    /* ~QPen(); */

    /* QPen &operator=(const QPen &pen) noexcept; */
    /* QPen(QPen &&other) noexcept */
    /*     : d(other.d) { other.d = nullptr; } */
    /* QPen &operator=(QPen &&other) noexcept */
    /* { qSwap(d, other.d); return *this; } */
    /* void swap(QPen &other) noexcept { qSwap(d, other.d); } */

    /* Qt::PenStyle style() const; */
    /* void setStyle(Qt::PenStyle); */

    /* QVector<qreal> dashPattern() const; */
    /* void setDashPattern(const QVector<qreal> &pattern); */

    /* qreal dashOffset() const; */
    /* void setDashOffset(qreal doffset); */

    /* qreal miterLimit() const; */
    /* void setMiterLimit(qreal limit); */

    /* qreal widthF() const; */
    /* void setWidthF(qreal width); */

    /* int width() const; */
    /* void setWidth(int width); */

    /* QColor color() const; */
    /* void setColor(const QColor &color); */

    /* QBrush brush() const; */
    /* void setBrush(const QBrush &brush); */

    /* bool isSolid() const; */

    /* Qt::PenCapStyle capStyle() const; */
    /* void setCapStyle(Qt::PenCapStyle pcs); */

    /* Qt::PenJoinStyle joinStyle() const; */
    /* void setJoinStyle(Qt::PenJoinStyle pcs); */

    /* bool isCosmetic() const; */
    /* void setCosmetic(bool cosmetic); */

    /* bool operator==(const QPen &p) const; */
    /* inline bool operator!=(const QPen &p) const { return !(operator==(p)); } */
    /* operator QVariant() const; */

    /* bool isDetached(); */
}; // QPen

class xQBrush : public QBrush
{
public:
    xQBrush();
/*     QBrush(Qt::BrushStyle bs); */
/*     QBrush(const QColor &color, Qt::BrushStyle bs=Qt::SolidPattern); */
/*     QBrush(Qt::GlobalColor color, Qt::BrushStyle bs=Qt::SolidPattern); */

/*     QBrush(const QColor &color, const QPixmap &pixmap); */
/*     QBrush(Qt::GlobalColor color, const QPixmap &pixmap); */
/*     QBrush(const QPixmap &pixmap); */
/*     QBrush(const QImage &image); */

/*     QBrush(const QBrush &brush); */

/*     QBrush(const QGradient &gradient); */

/*     ~QBrush(); */
/*     QBrush &operator=(const QBrush &brush); */
/*     inline QBrush &operator=(QBrush &&other) noexcept */
/*     { qSwap(d, other.d); return *this; } */
/*     inline void swap(QBrush &other) noexcept */
/*     { qSwap(d, other.d); } */

/*     operator QVariant() const; */

/*     inline Qt::BrushStyle style() const; */
/*     void setStyle(Qt::BrushStyle); */

/* #if QT_DEPRECATED_SINCE(5, 15) */
/*     QT_DEPRECATED_X("Use transform()") inline const QMatrix &matrix() const; */
/*     QT_DEPRECATED_X("Use setTransform()") void setMatrix(const QMatrix &mat); */
/* #endif // QT_DEPRECATED_SINCE(5, 15) */

/*     inline QTransform transform() const; */
/*     void setTransform(const QTransform &); */

/*     QPixmap texture() const; */
/*     void setTexture(const QPixmap &pixmap); */

/*     QImage textureImage() const; */
/*     void setTextureImage(const QImage &image); */

/*     inline const QColor &color() const; */
/*     void setColor(const QColor &color); */
/*     inline void setColor(Qt::GlobalColor color); */

/*     const QGradient *gradient() const; */

/*     bool isOpaque() const; */

/*     bool operator==(const QBrush &b) const; */
/*     inline bool operator!=(const QBrush &b) const { return !(operator==(b)); } */


/* public: */
/*     inline bool isDetached() const; */

}; // QBrush


class xQPixmap : public QPixmap
{
public:
    xQPixmap();
/*     explicit QPixmap(QPlatformPixmap *data); */
/*     QPixmap(int w, int h); */
/*     explicit QPixmap(const QSize &); */
/*     QPixmap(const QString& fileName, const char *format = nullptr, Qt::ImageConversionFlags flags = Qt::AutoColor); */
/* #ifndef QT_NO_IMAGEFORMAT_XPM */
/*     explicit QPixmap(const char * const xpm[]); */
/* #endif */
/*     QPixmap(const QPixmap &); */
/*     ~QPixmap(); */

/*     QPixmap &operator=(const QPixmap &); */
/*     inline QPixmap &operator=(QPixmap &&other) noexcept */
/*     { qSwap(data, other.data); return *this; } */
/*     inline void swap(QPixmap &other) noexcept */
/*     { qSwap(data, other.data); } */

/*     operator QVariant() const; */

/*     bool isNull() const; */
/*     int devType() const override; */

/*     int width() const; */
/*     int height() const; */
/*     QSize size() const; */
/*     QRect rect() const; */
/*     int depth() const; */

/*     static int defaultDepth(); */

/*     void fill(const QColor &fillColor = Qt::white); */
/* #if QT_DEPRECATED_SINCE(5, 13) */
/*     QT_DEPRECATED_X("Use QPainter or fill(QColor)") */
/*     void fill(const QPaintDevice *device, const QPoint &ofs); */
/*     QT_DEPRECATED_X("Use QPainter or fill(QColor)") */
/*     void fill(const QPaintDevice *device, int xofs, int yofs); */
/* #endif */

/*     QBitmap mask() const; */
/*     void setMask(const QBitmap &); */

/*     qreal devicePixelRatio() const; */
/*     void setDevicePixelRatio(qreal scaleFactor); */

/*     bool hasAlpha() const; */
/*     bool hasAlphaChannel() const; */

/* #ifndef QT_NO_IMAGE_HEURISTIC_MASK */
/*     QBitmap createHeuristicMask(bool clipTight = true) const; */
/* #endif */
/*     QBitmap createMaskFromColor(const QColor &maskColor, Qt::MaskMode mode = Qt::MaskInColor) const; */

/* #if QT_DEPRECATED_SINCE(5, 13) */
/*     QT_DEPRECATED_X("Use QScreen::grabWindow() instead") */
/*     static QPixmap grabWindow(WId, int x = 0, int y = 0, int w = -1, int h = -1); */
/*     QT_DEPRECATED_X("Use QWidget::grab() instead") */
/*     static QPixmap grabWidget(QObject *widget, const QRect &rect); */
/*     QT_DEPRECATED_X("Use QWidget::grab() instead") */
/*     static QPixmap grabWidget(QObject *widget, int x = 0, int y = 0, int w = -1, int h = -1); */
/* #endif */

/*     inline QPixmap scaled(int w, int h, Qt::AspectRatioMode aspectMode = Qt::IgnoreAspectRatio, */
/*                           Qt::TransformationMode mode = Qt::FastTransformation) const */
/*         { return scaled(QSize(w, h), aspectMode, mode); } */
/*     QPixmap scaled(const QSize &s, Qt::AspectRatioMode aspectMode = Qt::IgnoreAspectRatio, */
/*                    Qt::TransformationMode mode = Qt::FastTransformation) const; */
/*     QPixmap scaledToWidth(int w, Qt::TransformationMode mode = Qt::FastTransformation) const; */
/*     QPixmap scaledToHeight(int h, Qt::TransformationMode mode = Qt::FastTransformation) const; */
/* #if QT_DEPRECATED_SINCE(5, 15) */
/*     QT_DEPRECATED_X("Use transformed(const QTransform &, Qt::TransformationMode mode)") */
/*     QPixmap transformed(const QMatrix &, Qt::TransformationMode mode = Qt::FastTransformation) const; */
/*     QT_DEPRECATED_X("Use trueMatrix(const QTransform &m, int w, int h)") */
/*     static QMatrix trueMatrix(const QMatrix &m, int w, int h); */
/* #endif // QT_DEPRECATED_SINCE(5, 15) */
/*     QPixmap transformed(const QTransform &, Qt::TransformationMode mode = Qt::FastTransformation) const; */
/*     static QTransform trueMatrix(const QTransform &m, int w, int h); */

/*     QImage toImage() const; */
/*     static QPixmap fromImage(const QImage &image, Qt::ImageConversionFlags flags = Qt::AutoColor); */
/*     static QPixmap fromImageReader(QImageReader *imageReader, Qt::ImageConversionFlags flags = Qt::AutoColor); */
/*     static QPixmap fromImage(QImage &&image, Qt::ImageConversionFlags flags = Qt::AutoColor) */
/*     { */
/*         return fromImageInPlace(image, flags); */
/*     } */

/*     bool load(const QString& fileName, const char *format = nullptr, Qt::ImageConversionFlags flags = Qt::AutoColor); */
/*     bool loadFromData(const uchar *buf, uint len, const char* format = nullptr, Qt::ImageConversionFlags flags = Qt::AutoColor); */
/*     inline bool loadFromData(const QByteArray &data, const char* format = nullptr, Qt::ImageConversionFlags flags = Qt::AutoColor); */
/*     bool save(const QString& fileName, const char* format = nullptr, int quality = -1) const; */
/*     bool save(QIODevice* device, const char* format = nullptr, int quality = -1) const; */

/*     bool convertFromImage(const QImage &img, Qt::ImageConversionFlags flags = Qt::AutoColor); */

/*     inline QPixmap copy(int x, int y, int width, int height) const; */
/*     QPixmap copy(const QRect &rect = QRect()) const; */

/*     inline void scroll(int dx, int dy, int x, int y, int width, int height, QRegion *exposed = nullptr); */
/*     void scroll(int dx, int dy, const QRect &rect, QRegion *exposed = nullptr); */

/* #if QT_DEPRECATED_SINCE(5, 0) */
/*     QT_DEPRECATED inline int serialNumber() const { return cacheKey() >> 32; } */
/* #endif */
/*     qint64 cacheKey() const; */

/*     bool isDetached() const; */
/*     void detach(); */

/*     bool isQBitmap() const; */

/*     QPaintEngine *paintEngine() const override; */

/*     inline bool operator!() const { return isNull(); } */

/* #if QT_DEPRECATED_SINCE(5, 0) */
/*     QT_DEPRECATED inline QPixmap alphaChannel() const; */
/*     QT_DEPRECATED inline void setAlphaChannel(const QPixmap &); */
/* #endif */

/* protected: */
/*     int metric(PaintDeviceMetric) const override; */
/*     static QPixmap fromImageInPlace(QImage &image, Qt::ImageConversionFlags flags = Qt::AutoColor); */

/* public: */
/*     QPlatformPixmap* handle() const; */

}; // QPixmap

class xQIcon: public QIcon
{
public:
    xQIcon() noexcept;
    //    QIcon(const QPixmap &pixmap);
    //    QIcon(const QIcon &other);
/*     QIcon(QIcon &&other) noexcept */
/*         : d(other.d) */
/*     { other.d = nullptr; } */
/*     explicit QIcon(const QString &fileName); // file or resource name */
/*     explicit QIcon(QIconEngine *engine); */
/*     ~QIcon(); */
/*     QIcon &operator=(const QIcon &other); */
/*     inline QIcon &operator=(QIcon &&other) noexcept */
/*     { swap(other); return *this; } */
/*     inline void swap(QIcon &other) noexcept */
/*     { qSwap(d, other.d); } */

/*     operator QVariant() const; */

/*     QPixmap pixmap(const QSize &size, Mode mode = Normal, State state = Off) const; */
/*     inline QPixmap pixmap(int w, int h, Mode mode = Normal, State state = Off) const */
/*         { return pixmap(QSize(w, h), mode, state); } */
/*     inline QPixmap pixmap(int extent, Mode mode = Normal, State state = Off) const */
/*         { return pixmap(QSize(extent, extent), mode, state); } */
/*     QPixmap pixmap(QWindow *window, const QSize &size, Mode mode = Normal, State state = Off) const; */

/*     QSize actualSize(const QSize &size, Mode mode = Normal, State state = Off) const; */
/*     QSize actualSize(QWindow *window, const QSize &size, Mode mode = Normal, State state = Off) const; */

/*     QString name() const; */

/*     void paint(QPainter *painter, const QRect &rect, Qt::Alignment alignment = Qt::AlignCenter, Mode mode = Normal, State state = Off) const; */
/*     inline void paint(QPainter *painter, int x, int y, int w, int h, Qt::Alignment alignment = Qt::AlignCenter, Mode mode = Normal, State state = Off) const */
/*         { paint(painter, QRect(x, y, w, h), alignment, mode, state); } */

/*     bool isNull() const; */
/*     bool isDetached() const; */
/*     void detach(); */

/* #if QT_DEPRECATED_SINCE(5, 0) */
/*     QT_DEPRECATED inline int serialNumber() const { return cacheKey() >> 32; } */
/* #endif */
/*     qint64 cacheKey() const; */

/*     void addPixmap(const QPixmap &pixmap, Mode mode = Normal, State state = Off); */
/*     void addFile(const QString &fileName, const QSize &size = QSize(), Mode mode = Normal, State state = Off); */

/*     QList<QSize> availableSizes(Mode mode = Normal, State state = Off) const; */

/*     void setIsMask(bool isMask); */
/*     bool isMask() const; */

/*     static QIcon fromTheme(const QString &name); */
/*     static QIcon fromTheme(const QString &name, const QIcon &fallback); */
/*     static bool hasThemeIcon(const QString &name); */

/*     static QStringList themeSearchPaths(); */
/*     static void setThemeSearchPaths(const QStringList &searchpath); */

/*     static QStringList fallbackSearchPaths(); */
/*     static void setFallbackSearchPaths(const QStringList &paths); */

/*     static QString themeName(); */
/*     static void setThemeName(const QString &path); */

/*     static QString fallbackThemeName(); */
/*     static void setFallbackThemeName(const QString &name); */

/*     Q_DUMMY_COMPARISON_OPERATOR(QIcon) */

}; // QIcon

class xQImage : public QImage
{
public:
    xQImage() noexcept;
    xQImage(const QSize &size, Format format);
    xQImage(int width, int height, Format format);
/*     QImage(uchar *data, int width, int height, Format format, QImageCleanupFunction cleanupFunction = nullptr, void *cleanupInfo = nullptr); */
/*     QImage(const uchar *data, int width, int height, Format format, QImageCleanupFunction cleanupFunction = nullptr, void *cleanupInfo = nullptr); */
/*     QImage(uchar *data, int width, int height, int bytesPerLine, Format format, QImageCleanupFunction cleanupFunction = nullptr, void *cleanupInfo = nullptr); */
/*     QImage(const uchar *data, int width, int height, int bytesPerLine, Format format, QImageCleanupFunction cleanupFunction = nullptr, void *cleanupInfo = nullptr); */

/* #ifndef QT_NO_IMAGEFORMAT_XPM */
/*     explicit QImage(const char * const xpm[]); */
/* #endif */
/*     explicit QImage(const QString &fileName, const char *format = nullptr); */

/*     QImage(const QImage &); */
/*     inline QImage(QImage &&other) noexcept */
/*         : QPaintDevice(), d(nullptr) */
/*     { qSwap(d, other.d); } */
/*     ~QImage(); */

/*     QImage &operator=(const QImage &); */
/*     inline QImage &operator=(QImage &&other) noexcept */
/*     { qSwap(d, other.d); return *this; } */
/*     inline void swap(QImage &other) noexcept */
/*     { qSwap(d, other.d); } */

/*     bool isNull() const; */

/*     int devType() const override; */

/*     bool operator==(const QImage &) const; */
/*     bool operator!=(const QImage &) const; */
/*     operator QVariant() const; */
/*     void detach(); */
/*     bool isDetached() const; */

/*     QImage copy(const QRect &rect = QRect()) const; */
/*     inline QImage copy(int x, int y, int w, int h) const */
/*         { return copy(QRect(x, y, w, h)); } */

/*     Format format() const; */

/* #if defined(Q_COMPILER_REF_QUALIFIERS) && !defined(QT_COMPILING_QIMAGE_COMPAT_CPP) */
/*     Q_REQUIRED_RESULT Q_ALWAYS_INLINE QImage convertToFormat(Format f, Qt::ImageConversionFlags flags = Qt::AutoColor) const & */
/*     { return convertToFormat_helper(f, flags); } */
/*     Q_REQUIRED_RESULT Q_ALWAYS_INLINE QImage convertToFormat(Format f, Qt::ImageConversionFlags flags = Qt::AutoColor) && */
/*     { */
/*         if (convertToFormat_inplace(f, flags)) */
/*             return std::move(*this); */
/*         else */
/*             return convertToFormat_helper(f, flags); */
/*     } */
/* #else */
/*     Q_REQUIRED_RESULT QImage convertToFormat(Format f, Qt::ImageConversionFlags flags = Qt::AutoColor) const; */
/* #endif */
/*     Q_REQUIRED_RESULT QImage convertToFormat(Format f, const QVector<QRgb> &colorTable, Qt::ImageConversionFlags flags = Qt::AutoColor) const; */
/*     bool reinterpretAsFormat(Format f); */

/*     void convertTo(Format f, Qt::ImageConversionFlags flags = Qt::AutoColor); */

/*     int width() const; */
/*     int height() const; */
/*     QSize size() const; */
/*     QRect rect() const; */

/*     int depth() const; */
/*     int colorCount() const; */
/*     int bitPlaneCount() const; */

/*     QRgb color(int i) const; */
/*     void setColor(int i, QRgb c); */
/*     void setColorCount(int); */

/*     bool allGray() const; */
/*     bool isGrayscale() const; */

/*     uchar *bits(); */
/*     const uchar *bits() const; */
/*     const uchar *constBits() const; */

/* #if QT_DEPRECATED_SINCE(5, 10) */
/*     QT_DEPRECATED_X("Use sizeInBytes") int byteCount() const; */
/* #endif */
/*     qsizetype sizeInBytes() const; */

/*     uchar *scanLine(int); */
/*     const uchar *scanLine(int) const; */
/*     const uchar *constScanLine(int) const; */
/* #if QT_VERSION >= QT_VERSION_CHECK(6,0,0) */
/*     qsizetype bytesPerLine() const; */
/* #else */
/*     int bytesPerLine() const; */
/* #endif */

/*     bool valid(int x, int y) const; */
/*     bool valid(const QPoint &pt) const; */

/*     int pixelIndex(int x, int y) const; */
/*     int pixelIndex(const QPoint &pt) const; */

/*     QRgb pixel(int x, int y) const; */
/*     QRgb pixel(const QPoint &pt) const; */

/*     void setPixel(int x, int y, uint index_or_rgb); */
/*     void setPixel(const QPoint &pt, uint index_or_rgb); */

/*     QColor pixelColor(int x, int y) const; */
/*     QColor pixelColor(const QPoint &pt) const; */

/*     void setPixelColor(int x, int y, const QColor &c); */
/*     void setPixelColor(const QPoint &pt, const QColor &c); */

/*     QVector<QRgb> colorTable() const; */
/* #if QT_VERSION >= QT_VERSION_CHECK(6,0,0) */
/*     void setColorTable(const QVector<QRgb> &colors); */
/* #else */
/*     void setColorTable(const QVector<QRgb> colors); */
/* #endif */

/*     qreal devicePixelRatio() const; */
/*     void setDevicePixelRatio(qreal scaleFactor); */

/*     void fill(uint pixel); */
/*     void fill(const QColor &color); */
/*     void fill(Qt::GlobalColor color); */


/*     bool hasAlphaChannel() const; */
/*     void setAlphaChannel(const QImage &alphaChannel); */
/* #if QT_DEPRECATED_SINCE(5, 15) */
/*     QT_DEPRECATED_X("Use convertToFormat(QImage::Format_Alpha8)") */
/*     QImage alphaChannel() const; */
/* #endif */
/*     QImage createAlphaMask(Qt::ImageConversionFlags flags = Qt::AutoColor) const; */
/* #ifndef QT_NO_IMAGE_HEURISTIC_MASK */
/*     QImage createHeuristicMask(bool clipTight = true) const; */
/* #endif */
/*     QImage createMaskFromColor(QRgb color, Qt::MaskMode mode = Qt::MaskInColor) const; */

/*     inline QImage scaled(int w, int h, Qt::AspectRatioMode aspectMode = Qt::IgnoreAspectRatio, */
/*                         Qt::TransformationMode mode = Qt::FastTransformation) const */
/*         { return scaled(QSize(w, h), aspectMode, mode); } */
/*     QImage scaled(const QSize &s, Qt::AspectRatioMode aspectMode = Qt::IgnoreAspectRatio, */
/*                  Qt::TransformationMode mode = Qt::FastTransformation) const; */
/*     QImage scaledToWidth(int w, Qt::TransformationMode mode = Qt::FastTransformation) const; */
/*     QImage scaledToHeight(int h, Qt::TransformationMode mode = Qt::FastTransformation) const; */
/* #if QT_DEPRECATED_SINCE(5, 15) */
/*     QT_DEPRECATED_X("Use transformed(const QTransform &matrix, Qt::TransformationMode mode)") */
/*     QImage transformed(const QMatrix &matrix, Qt::TransformationMode mode = Qt::FastTransformation) const; */
/*     QT_DEPRECATED_X("trueMatrix(const QTransform &, int w, int h)") */
/*     static QMatrix trueMatrix(const QMatrix &, int w, int h); */
/* #endif // QT_DEPRECATED_SINCE(5, 15) */
/*     QImage transformed(const QTransform &matrix, Qt::TransformationMode mode = Qt::FastTransformation) const; */
/*     static QTransform trueMatrix(const QTransform &, int w, int h); */
/* #if defined(Q_COMPILER_REF_QUALIFIERS) && !defined(QT_COMPILING_QIMAGE_COMPAT_CPP) */
/*     QImage mirrored(bool horizontally = false, bool vertically = true) const & */
/*         { return mirrored_helper(horizontally, vertically); } */
/*     QImage &&mirrored(bool horizontally = false, bool vertically = true) && */
/*         { mirrored_inplace(horizontally, vertically); return std::move(*this); } */
/*     QImage rgbSwapped() const & */
/*         { return rgbSwapped_helper(); } */
/*     QImage &&rgbSwapped() && */
/*         { rgbSwapped_inplace(); return std::move(*this); } */
/* #else */
/*     QImage mirrored(bool horizontally = false, bool vertically = true) const; */
/*     QImage rgbSwapped() const; */
/* #endif */
/*     void invertPixels(InvertMode = InvertRgb); */

/*     QColorSpace colorSpace() const; */
/*     QImage convertedToColorSpace(const QColorSpace &) const; */
/*     void convertToColorSpace(const QColorSpace &); */
/*     void setColorSpace(const QColorSpace &); */

/*     void applyColorTransform(const QColorTransform &transform); */

/*     bool load(QIODevice *device, const char* format); */
/*     bool load(const QString &fileName, const char *format = nullptr); */
/*     bool loadFromData(const uchar *buf, int len, const char *format = nullptr); */
/*     inline bool loadFromData(const QByteArray &data, const char *aformat = nullptr) */
/*         { return loadFromData(reinterpret_cast<const uchar *>(data.constData()), data.size(), aformat); } */

/*     bool save(const QString &fileName, const char *format = nullptr, int quality = -1) const; */
/*     bool save(QIODevice *device, const char *format = nullptr, int quality = -1) const; */

/*     static QImage fromData(const uchar *data, int size, const char *format = nullptr); */
/*     inline static QImage fromData(const QByteArray &data, const char *format = nullptr) */
/*         { return fromData(reinterpret_cast<const uchar *>(data.constData()), data.size(), format); } */

/* #if QT_DEPRECATED_SINCE(5, 0) */
/*     QT_DEPRECATED inline int serialNumber() const { return cacheKey() >> 32; } */
/* #endif */
/*     qint64 cacheKey() const; */

/*     QPaintEngine *paintEngine() const override; */

/*     // Auxiliary data */
/*     int dotsPerMeterX() const; */
/*     int dotsPerMeterY() const; */
/*     void setDotsPerMeterX(int); */
/*     void setDotsPerMeterY(int); */
/*     QPoint offset() const; */
/*     void setOffset(const QPoint&); */

/*     QStringList textKeys() const; */
/*     QString text(const QString &key = QString()) const; */
/*     void setText(const QString &key, const QString &value); */

/*     QPixelFormat pixelFormat() const noexcept; */
/*     static QPixelFormat toPixelFormat(QImage::Format format) noexcept; */
/*     static QImage::Format toImageFormat(QPixelFormat format) noexcept; */

/*     // Platform specific conversion functions */
/* #if defined(Q_OS_DARWIN) || defined(Q_QDOC) */
/*     CGImageRef toCGImage() const Q_DECL_CF_RETURNS_RETAINED; */
/* #endif */

/* #if QT_DEPRECATED_SINCE(5, 0) */
/*     QT_DEPRECATED inline QString text(const char *key, const char *lang = nullptr) const; */
/*     QT_DEPRECATED inline QList<QImageTextKeyLang> textList() const; */
/*     QT_DEPRECATED inline QStringList textLanguages() const; */
/*     QT_DEPRECATED inline QString text(const QImageTextKeyLang&) const; */
/*     QT_DEPRECATED inline void setText(const char* key, const char* lang, const QString&); */
/* #endif */

/* #if QT_DEPRECATED_SINCE(5, 0) */
/*     QT_DEPRECATED inline int numColors() const; */
/*     QT_DEPRECATED inline void setNumColors(int); */
/*     QT_DEPRECATED inline int numBytes() const; */
/* #endif */

/* protected: */
/*     virtual int metric(PaintDeviceMetric metric) const override; */
/*     QImage mirrored_helper(bool horizontal, bool vertical) const; */
/*     QImage rgbSwapped_helper() const; */
/*     void mirrored_inplace(bool horizontal, bool vertical); */
/*     void rgbSwapped_inplace(); */
/*     QImage convertToFormat_helper(Format format, Qt::ImageConversionFlags flags) const; */
/*     bool convertToFormat_inplace(Format format, Qt::ImageConversionFlags flags); */
/*     QImage smoothScaled(int w, int h) const; */

}; // QImage



class xQClipboard : public QClipboard
{
public:
    void clear(Mode mode = Clipboard);

    bool supportsSelection() const;
    bool supportsFindBuffer() const;

    bool ownsSelection() const;
    bool ownsClipboard() const;
    bool ownsFindBuffer() const;

/*     QString text(Mode mode = Clipboard) const; */
/*     QString text(QString& subtype, Mode mode = Clipboard) const; */
/*     void setText(const QString &, Mode mode = Clipboard); */

/*     const QMimeData *mimeData(Mode mode = Clipboard ) const; */
/*     void setMimeData(QMimeData *data, Mode mode = Clipboard); */

/*     QImage image(Mode mode = Clipboard) const; */
/*     QPixmap pixmap(Mode mode = Clipboard) const; */
/*     void setImage(const QImage &, Mode mode  = Clipboard); */
/*     void setPixmap(const QPixmap &, Mode mode  = Clipboard); */

/* Q_SIGNALS: */
/*     void changed(QClipboard::Mode mode); */
/*     void selectionChanged(); */
/*     void findBufferChanged(); */
/*     void dataChanged(); */

}; // QClipboard


class xQWindow : public QWindow
{
 public:

};

class xQSurface : public QSurface {
 public:
};

class xQDesktopServices : public QDesktopServices
{
public:
    static bool openUrl(const QUrl &url);
    static void setUrlHandler(const QString &scheme, QObject *receiver, const char *method);
    static void unsetUrlHandler(const QString &scheme);

/* #if QT_DEPRECATED_SINCE(5, 0) */
/*     //Must match QStandardPaths::StandardLocation */
/*     enum StandardLocation { */
/*         DesktopLocation, */
/*         DocumentsLocation, */
/*         FontsLocation, */
/*         ApplicationsLocation, */
/*         MusicLocation, */
/*         MoviesLocation, */
/*         PicturesLocation, */
/*         TempLocation, */
/*         HomeLocation, */
/*         DataLocation, */
/*         CacheLocation */
/*     }; */

/*     QT_DEPRECATED static QString storageLocation(StandardLocation type) { */
/*         return storageLocationImpl(static_cast<QStandardPaths::StandardLocation>(type)); */
/*     } */
/*     QT_DEPRECATED static QString displayName(StandardLocation type) { */
/*         return QStandardPaths::displayName(static_cast<QStandardPaths::StandardLocation>(type)); */
/*     } */
/* #endif */
}; // QDesktopServices


class xQGuiApplication : public QGuiApplication {
 public:
};


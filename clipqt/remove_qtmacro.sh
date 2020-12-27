set -x

sed -i 's/Q_REQUIRED_RESULT//g' *0.h
sed -i 's/Q_DECL_CONSTEXPR//g' *0.h
sed -i 's/Q_DECL_RELAXED_CONSTEXPR//g' *0.h
sed -i 's/Q_CORE_EXPORT//g' *0.h
sed -i 's/Q_GUI_EXPORT//g' *0.h
sed -i 's/Q_WIDGETS_EXPORT//g' *0.h

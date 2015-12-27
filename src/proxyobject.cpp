
#include <QtCore>
#include <QtGui>
#include <QtWidgets>

class NPushButton : public QPushButton
{
    Q_OBJECT;
public:

};

// QtClass_SlotProxy
class SlotProxy_Base : public QObject
{
public:
    // 声明slot函数指针，函数指针中的类型是rust能接收到的。
    SlotProxy_Base() : QObject() {}
    // 定义所有可能slot方法，
};


class NPushButton_SlotProxy : public QObject
{
    Q_OBJECT;
public:
    NPushButton_SlotProxy():QObject(){}


    // 函数指针成员变量
    void (*slot_func_clicked_mangled)(bool) = NULL;
    void (*slot_func_pressed_mangled)()  = NULL;
    void (*slot_func_released_mangled)() = NULL;

    // slot方法
    void slot_proxy_func_clicked(bool checked = false);
    void slot_proxy_func_pressed();
    void slot_proxy_func_released();
};
#include "proxyobject.moc"

void NPushButton_SlotProxy::slot_proxy_func_clicked(bool checked)
{
    if (this->slot_func_clicked_mangled != NULL) {
        // do smth...
        // 在这转换类型吗？
        this->slot_func_clicked_mangled(checked);
    }
}
void NPushButton_SlotProxy::slot_proxy_func_pressed()
{
    if (this->slot_func_pressed_mangled != NULL) {
        this->slot_func_pressed_mangled();
    }
}
void NPushButton_SlotProxy::slot_proxy_func_released()
{
    if (this->slot_func_released_mangled != NULL) {
        this->slot_func_released_mangled();
    }
}


extern "C" {
    NPushButton_SlotProxy *NPushButton_SlotProxy_new()
    {
        return new NPushButton_SlotProxy();
    }

    void NPushButton_SlotProxy_delete(NPushButton_SlotProxy *that)
    {
        delete that;
    }

    void QAbstractButton_SlotProxy_connect_pressed_mangled
    (
     QObject *sender, void *fptr
     // QAbstractButton_SlotProxy *that,
     // decltype(that->slot_func_pressed_mangled) fptr
     ) {
        NPushButton_SlotProxy *that = new NPushButton_SlotProxy();
        that->slot_func_pressed_mangled = (decltype(that->slot_func_pressed_mangled))fptr;
        QObject::connect((NPushButton*)sender, &QPushButton::pressed,
                         that, &NPushButton_SlotProxy::slot_proxy_func_pressed);
    }

    void QAbstractButton_SlotProxy_disconnect_pressed_mangled(NPushButton_SlotProxy *that)
    {
        that->disconnect();
        delete that;
    }
};


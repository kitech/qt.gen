
#include <QtCore>
#include <QtGui>
#include <QtWidgets>

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
    void (*slot_func_clicked)(bool) = NULL;
    void (*slot_func_pressed)()  = NULL;
    void (*slot_func_released)() = NULL;

    // slot方法
    void slot_proxy_func_clicked(bool checked = false);
    void slot_proxy_func_pressed();
    void slot_proxy_func_released();
};
#include "proxyobject.moc"

void NPushButton_SlotProxy::slot_proxy_func_clicked(bool checked)
{
    if (this->slot_func_clicked != NULL) {
        // do smth...
        // 在这转换类型吗？
        this->slot_func_clicked(checked);
    }
}
void NPushButton_SlotProxy::slot_proxy_func_pressed()
{
    if (this->slot_func_pressed != NULL) {
        this->slot_func_pressed();
    }
}
void NPushButton_SlotProxy::slot_proxy_func_released()
{
    if (this->slot_func_released != NULL) {
        this->slot_func_released();
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

};


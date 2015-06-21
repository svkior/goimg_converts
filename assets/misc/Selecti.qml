import QtQuick 2.3

Rectangle {
    id: container
    property alias fontSize : defText.font.pixelSize
    property string btnColor : "#003300"
    property string btnColorSelected : "#006600"
    property string btnText: fileName
    property string btnMode : ""
    property string btnSelected : "1"
    property bool canProvideClick // Нужно для того чтобы нажатие пошло дальше

    signal clicked()
    color: btnTrueColor

    property bool isEntered : false
    property bool isClicked: false

    property string btnClickedColor : (isClicked ? Qt.darker(btnPreColor, 2) : btnPreColor)
    property string btnPreColor : (isEntered ? Qt.lighter(btnColor, 1.4) : btnColor)
    property string btnTrueColor : (btnMode == btnSelected ? btnColorSelected : btnClickedColor)


    border.color: "#a37777" // Поставил коментарий
    border.width: 2
    radius: 5

    width: 200
    height: 40

    DefText {
        id:defText
        anchors.fill: parent
        horizontalAlignment: Text.AlignHCenter
        verticalAlignment: Text.AlignVCenter
        font.pixelSize: 24
        color: container.border.color
        text: container.btnText
    }
    gradient: Gradient {
        GradientStop {
            position: 0.0; color: (btnMode == btnSelected ? "#202020" : Qt.darker(btnClickedColor))
        }
        GradientStop {
            position: 0.3
            color:
                    btnTrueColor
        }
        GradientStop {
            position: 1.0; color: (btnMode == btnSelected ? "#202020" : Qt.darker(btnClickedColor))
        }
    }
    MouseArea {
        id: mouseArea
        anchors.fill: parent
        hoverEnabled: true
        onEntered: container.isEntered = true
        onExited:   {
            container.isEntered = false
            container.isClicked = false
            container.canProvideClick = false // Если вышли при нажатой мыши, то не вызывать действие
        }
        onPressed: {
            container.canProvideClick = true;
            container.isClicked = true
        }
        onReleased: {
            container.isClicked = false
            if(container.canProvideClick){
                container.clicked(container.btnText)
            } else {
            }
        }
    }
}
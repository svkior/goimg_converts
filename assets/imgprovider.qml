import QtQuick 2.2
import QtQuick.Controls 1.2
import QtQuick.Dialogs 1.1
import QtQuick.Layouts 1.1
import QtQuick.Window 2.0
import Qt.labs.folderlistmodel 2.1
import "./misc"

Item {
	width: 1024
	height: 768
	SystemPalette {id: palette}
	clip: true

	FileDialog {
		id:fileDialog
		visible: false
		modality: Qt.NonModal
		title: "Выберите папку с изображениями"
		selectExisting: true
		selectMultiple: false
		selectFolder: true
		sidebarVisible: true
		onAccepted: {
			console.log("Accepted: ", fileUrls)
			folderName.text = fileUrls[0]
			folderModel.folder = fileUrls[0]
		}
		onRejected: {
			console.log("Rejected")
		}
	}	


	Rectangle{
		id: bottomBar
		color: "#000000"
		anchors{
			top: parent.top
			left: parent.left
			right: parent.right
			bottom: parent.bottom
			} 

		Column {
			spacing: 6
			anchors.top: parent.top
			anchors.bottom: parent.bottom
			anchors.left: parent.left
			anchors.leftMargin: 12

			Row {
				Button {
					text: "Open"
					onClicked: fileDialog.open()
				}
				TextEdit {
					id: folderName
					width: 240
					text: ""
					font.family: "Helvetica"
					font.pointSize: 20
					color: "#0000dd"
					focus: true
				}

			}

			Row {
				spacing: 5
				ListView {
					id: filesList
				   width: 200
	 			   height: parent.height
	 			   spacing: 7

	    			FolderListModel {
	        		id: folderModel
	        		showDirs: false
	        		nameFilters: ["*.jpg", "*.jpeg", "*.png"]
	    			}

	    			

	    			model: folderModel
	    			delegate: Selecti{
	    				btnMode: "image://pwd/" + fileURL
	    				btnSelected: imgSrc.source
	    				onClicked: {
	    					//console.log("Clicked:", fileURL);
	    					imgSrc.source = "image://pwd/" + fileURL
	    				}
	    			}
				}
				Image {
					id: imgSrc
					width: 800
					height: 600
					fillMode: Image.PreserveAspectFit
			    	source: "./TTS-watermark-white.svg"
				}	


			}


		}
	}

}


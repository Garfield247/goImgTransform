package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MyMainWindow struct {
	*walk.MainWindow
	edit *walk.TextEdit
}

func main() {
	mw := &MyMainWindow{}
	err := MainWindow{
		AssignTo: &mw.MainWindow, //窗口重定向至mw，重定向后可由重定向变量控制控件
		// Icon:     "test.ico",     //窗体图标
		Title:   "文件选择对话框", //标题
		MinSize: Size{Width: 150, Height: 200},
		Size:    Size{300, 400},
		Layout:  VBox{}, //样式，纵向
		Children: []Widget{ //控件组
			TextEdit{
				AssignTo: &mw.edit,
			},
			PushButton{
				Text:      "打开",
				OnClicked: mw.selectFile, //点击事件响应函数
			},
			PushButton{
				Text:      "另存为",
				OnClicked: mw.saveFile,
			},
		},
	}.Create() //创建

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	mw.Run() //运行
}

func (mw *MyMainWindow) selectFile() {

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	dlg.Filter = "可执行文件 (*.exe)|*.exe|所有文件 (*.*)|*.*"

	mw.edit.SetText("") //通过重定向变量设置TextEdit的Text
	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.edit.AppendText("Error : File Open\r\n")
		return
	} else if !ok {
		mw.edit.AppendText("Cancel\r\n")
		return
	}
	s := fmt.Sprintf("Select : %s\r\n", dlg.FilePath)
	mw.edit.AppendText(s)
}

func (mw *MyMainWindow) saveFile() {

	dlg := new(walk.FileDialog)
	dlg.Title = "另存为"

	if ok, err := dlg.ShowSave(mw); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	} else if !ok {
		fmt.Fprintln(os.Stderr, "Cancel")
		return
	}

	data := mw.edit.Text()
	filename := dlg.FilePath
	f, err := os.Open(filename)
	if err != nil {
		f, _ = os.Create(filename)
	} else {
		f.Close()
		//打开文件，参数：文件路径及名称，打开方式，控制权限
		f, err = os.OpenFile(filename, os.O_WRONLY, 0x666)
	}
	if len(data) == 0 {
		f.Close()
		return
	}
	io.Copy(f, strings.NewReader(data))
	f.Close()
}

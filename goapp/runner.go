package main

import (
	"fmt"
	"github.com/nvsoft/cef"
	//"github.com/nvsoft/win"
	"github.com/nvsoft/goapp/navigator"
	"os"
	"strconv"
	"strings"
	//"syscall"
	"time"
)

func working() {
	for {
		browser, b := cef.MainBrowser()
		if b {
			fmt.Printf("browser id=%v pid:%v\n", browser.Id, os.Getpid())
			doIt(browser)
		}
		time.Sleep(3 * time.Second)
	}
}

func doIt(browser *cef.Browser) {
	// 1.加载百度首页
	browser.LoadURL("http://www.baidu.com")
	time.Sleep(5 * time.Second)

	url := browser.GetURL()
	fmt.Printf("Url=%v\n", url)

	src := browser.GetSource()

	fmt.Printf("src=%v\n", len(src))

	//js := `function() { return "1"; }();`
	//browser.Eval(js)

	// 执行Js
	//browser.ExecuteJavaScript(`app.cefResult("a");`, "", 1)
	result := browser.ExecuteJavaScriptWithResult(`(function() { app.cefResult("13"); })();`)
	fmt.Printf("Eval Js. result=%v\n", result)

	//browser.InjectJs("js/jquery.min.js")

	// http://stackoverflow.com/questions/12605315/html5-getimagedata-without-canvas
	// http://blog.csdn.net/hursing/article/details/12868109

	//browser.InjectJs("js/imageutil.js")
	//browser.InjectJs("js/html2canvas.js")
	/*filename := browser.ExecuteJavaScriptWithResult(`
	(function() {
		var img = $("#lg img")[0];
		var imageData = captureImage(img);
		var filename = app.renderImage(imageData.data, imageData.width, imageData.height);
		return filename;
	})();`)

	fmt.Printf("filename: %v\n", filename)
	*/

	/*
		filename := browser.ExecuteJavaScriptWithResult(`
		html2canvas(document.body, {
		  onrendered: function(canvas) {
			var imageData = canvasToData(canvas);
			var filename = app.renderImage(imageData.data, imageData.width, imageData.height);
			app.cefResult(filename);
		  }
		});
		`)
		fmt.Printf("filename: %v\n", filename)
	*/

	/*c := `var canvas = document.createElement('canvas');
		var context = canvas.getContext('2d');
		var img = $("#lg img")[0];
		//var img = document.getElementById('tulip');
		canvas.width = img.width;
		canvas.height = img.height;
		context.drawImage(img, 0, 0);//, img.width, img.height
		var imageData = context.getImageData(0, 0, img.width, img.height);
		//alert('3:' + img.width + '/' + img.height + '/' + imageData.data.length);
		var txtFile = '';

		//console.log('image.data.length=' + imageData.data.length);

		var dataArray = new Array(imageData.data.length);
	    for (var i = 0; i < dataArray.length; i++) {
	        dataArray[i] = imageData.data[i];
	    }
	    var strImageData = dataArray.toString();

		//console.log("strImageData=" + strImageData);

		var filename = app.set_image_data(strImageData, img.width, img.height);
		alert(filename);
		//alert('4');`
		//c := `var canvas = $("#lg img");var context = canvas[0].getContext("2d");alert(context);`
	*/

	// 获取图片
	//browser.ExecuteJavaScript(c, "", 1)

	// 发送进程消息
	//browser.SendProcessMessageTest()

	// 模拟点击
	code := `
    function getOffset( el ) {
        var _x = 0;
        var _y = 0;
        while( el && !isNaN( el.offsetLeft ) && !isNaN( el.offsetTop ) ) {
            _x += el.offsetLeft - el.scrollLeft;
            _y += el.offsetTop - el.scrollTop;
            // chrome/safari
            //if ($.browser.webkit) {
                el = el.parentNode;
            //} else {
                // firefox/IE
                //el = el.offsetParent;
            //}
        }
        return { left: _x, top: _y };
    }
    var e = document.querySelector("input[name=wd]");
    var offset = getOffset(e);
    app.cefResult(offset.left + "," + offset.top);
    `
	strOffset := browser.ExecuteJavaScriptWithResult(code)
	fmt.Printf("strOffset:%v\n", strOffset)

	ss := strings.Split(strOffset, ",")

	//h := browser.GetHost()

	if len(ss) == 2 {
		x, _ := strconv.Atoi(ss[0])
		y, _ := strconv.Atoi(ss[1])
		// 执行点击
		fmt.Printf("x=%v,y=%v\n", x, y)
		navigator.InjectMouseClick(browser, x+4, y+4)
	}

	fmt.Printf("点击完毕\n")
	time.Sleep(3 * time.Second)

	//h.SetFocus(true)

	// 模拟输入
	fmt.Printf("模拟输入\n")
	var text = "Abcdd@ffaff"
	navigator.InjectKeyPress(browser, text)

	time.Sleep(10 * time.Second)
}

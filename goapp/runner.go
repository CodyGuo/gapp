package main

import (
	"fmt"
	"github.com/nvsoft/cef"
	"os"
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

	/*
		var c=document.getElementById("myCanvas");
		var ctx=c.getContext("2d");
	*/
	//browser.InjectJs("js/jquery.min.js")

	// http://stackoverflow.com/questions/12605315/html5-getimagedata-without-canvas
	// http://blog.csdn.net/hursing/article/details/12868109

	/*c := `console.log("$=" + $);
		var canvas = $("#lg img")[0];
		var context = canvas.getContext("2d");
		var imageData = context.getImageData(0,0,500,500);
		var data = imageData.data;
	   app.set_image_data(data,500,500);`*/
	//browser.InjectJs("js/jquery.min.js")
	browser.InjectJs("js/imageutil.js")
	browser.InjectJs("js/html2canvas.js")
	/*filename := browser.ExecuteJavaScriptWithResult(`
	(function() {
		var img = $("#lg img")[0];
		var imageData = captureImage(img);
		var filename = app.renderImage(imageData.data, imageData.width, imageData.height);
		return filename;
	})();`)

	fmt.Printf("filename: %v\n", filename)
	*/

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

	time.Sleep(10 * time.Second)
}

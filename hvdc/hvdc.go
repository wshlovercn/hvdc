package main

import "hvdc/baselib/app"

func main()  {
	appDelegate := &HdvcServer{}
	application := app.NewApplication(appDelegate)
	application.Run()
}

package main

import (
	"context"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/opentracing/opentracing-go/log"
	"shop/fronted/web/controllers"
	"shop/fronted/middleware"
	"shop/common"
	"shop/repositories"
	"shop/services"
	"time"
)

func main() {
	// 1. 创建iris实例
	app := iris.New()
	// 2. 设置错误模式,在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	// 3.注册模板
	tmplate := iris.HTML("./fronted/web/views", "html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	// 4. 设置模板目标
	app.HandleDir("/public", "./fronted/web/public")
	// 出现异常跳转至指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错!"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	// 连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sess := sessions.New(sessions.Config{
		Cookie:  "helloworld",
		Expires: 60 * time.Minute,
	})
	// 5. 注册控制器
	userRepository := repositories.NewUserRepository("user", db)
	userService := services.NewService(userRepository)
	user := mvc.New(app.Party("/user"))
	user.Register(userService, ctx, sess.Start)
	user.Handle(new(controllers.UserController))

	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	orderRepository := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)

	productParty := app.Party("/product")
	product := mvc.New(productParty)
	productParty.Use(middleware.AuthConProduct)
	product.Register(productService, orderService, sess.Start)
	product.Handle(new(controllers.ProductController))

	// 6. 启动服务
	app.Run(iris.Addr("0.0.0.0:8082"),
		//iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}

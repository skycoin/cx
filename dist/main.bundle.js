webpackJsonp(["main"],{

/***/ "../../../../../src/$$_gendir lazy recursive":
/***/ (function(module, exports) {

function webpackEmptyAsyncContext(req) {
	// Here Promise.resolve().then() is used instead of new Promise() to prevent
	// uncatched exception popping up in devtools
	return Promise.resolve().then(function() {
		throw new Error("Cannot find module '" + req + "'.");
	});
}
webpackEmptyAsyncContext.keys = function() { return []; };
webpackEmptyAsyncContext.resolve = webpackEmptyAsyncContext;
module.exports = webpackEmptyAsyncContext;
webpackEmptyAsyncContext.id = "../../../../../src/$$_gendir lazy recursive";

/***/ }),

/***/ "../../../../../src/app/app.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/app.component.html":
/***/ (function(module, exports) {

module.exports = "\n\n<app-header></app-header>\n\n<div class =\"container\">\n<router-outlet></router-outlet>\n</div>\n\n"

/***/ }),

/***/ "../../../../../src/app/app.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AppComponent; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__ = __webpack_require__("../../../platform-browser/@angular/platform-browser.es5.js");
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var AppComponent = (function () {
    function AppComponent(titleService) {
        this.titleService = titleService;
    }
    return AppComponent;
}());
AppComponent = __decorate([
    Object(__WEBPACK_IMPORTED_MODULE_0__angular_core__["o" /* Component */])({
        selector: 'app-root',
        template: __webpack_require__("../../../../../src/app/app.component.html"),
        styles: [__webpack_require__("../../../../../src/app/app.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["b" /* Title */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["b" /* Title */]) === "function" && _a || Object])
], AppComponent);

var _a;
//# sourceMappingURL=app.component.js.map

/***/ }),

/***/ "../../../../../src/app/app.module.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AppModule; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__ = __webpack_require__("../../../platform-browser/@angular/platform-browser.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_forms__ = __webpack_require__("../../../forms/@angular/forms.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__angular_http__ = __webpack_require__("../../../http/@angular/http.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__app_component__ = __webpack_require__("../../../../../src/app/app.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6__components_header_header_comonent__ = __webpack_require__("../../../../../src/app/components/header/header.comonent.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7__components_home_home_comonent__ = __webpack_require__("../../../../../src/app/components/home/home.comonent.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8__components_tutorial_tutorial_comonent__ = __webpack_require__("../../../../../src/app/components/tutorial/tutorial.comonent.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9__components_examples_examples_comonent__ = __webpack_require__("../../../../../src/app/components/examples/examples.comonent.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10__components_about_about_comonent__ = __webpack_require__("../../../../../src/app/components/about/about.comonent.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11__components_faq_faq_comonent__ = __webpack_require__("../../../../../src/app/components/faq/faq.comonent.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12__app_routing__ = __webpack_require__("../../../../../src/app/app.routing.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13_ngx_bootstrap_dropdown__ = __webpack_require__("../../../../ngx-bootstrap/dropdown/index.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14_ngx_bootstrap_tooltip__ = __webpack_require__("../../../../ngx-bootstrap/tooltip/index.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_ngx_bootstrap_modal__ = __webpack_require__("../../../../ngx-bootstrap/modal/index.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_16_ngx_bootstrap_collapse__ = __webpack_require__("../../../../ngx-bootstrap/collapse/index.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_17__services_api_service__ = __webpack_require__("../../../../../src/app/services/api.service.ts");
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};


















var AppModule = (function () {
    function AppModule() {
    }
    return AppModule;
}());
AppModule = __decorate([
    Object(__WEBPACK_IMPORTED_MODULE_1__angular_core__["M" /* NgModule */])({
        declarations: [
            __WEBPACK_IMPORTED_MODULE_5__app_component__["a" /* AppComponent */],
            __WEBPACK_IMPORTED_MODULE_6__components_header_header_comonent__["a" /* HeaderComponent */],
            __WEBPACK_IMPORTED_MODULE_7__components_home_home_comonent__["a" /* HomeComponent */],
            __WEBPACK_IMPORTED_MODULE_8__components_tutorial_tutorial_comonent__["a" /* TutorialComponent */],
            __WEBPACK_IMPORTED_MODULE_9__components_examples_examples_comonent__["a" /* ExamplesComponent */],
            __WEBPACK_IMPORTED_MODULE_10__components_about_about_comonent__["a" /* AboutComponent */],
            __WEBPACK_IMPORTED_MODULE_11__components_faq_faq_comonent__["a" /* FAQComponent */]
        ],
        imports: [
            __WEBPACK_IMPORTED_MODULE_4__angular_http__["c" /* HttpModule */],
            __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__["a" /* BrowserModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_forms__["a" /* FormsModule */],
            __WEBPACK_IMPORTED_MODULE_3__angular_router__["b" /* RouterModule */].forRoot(__WEBPACK_IMPORTED_MODULE_12__app_routing__["a" /* AppRoutes */]),
            __WEBPACK_IMPORTED_MODULE_13_ngx_bootstrap_dropdown__["a" /* BsDropdownModule */].forRoot(),
            __WEBPACK_IMPORTED_MODULE_14_ngx_bootstrap_tooltip__["a" /* TooltipModule */].forRoot(),
            __WEBPACK_IMPORTED_MODULE_15_ngx_bootstrap_modal__["a" /* ModalModule */].forRoot(),
            __WEBPACK_IMPORTED_MODULE_16_ngx_bootstrap_collapse__["a" /* CollapseModule */].forRoot()
        ],
        providers: [__WEBPACK_IMPORTED_MODULE_17__services_api_service__["a" /* ApiService */]],
        bootstrap: [__WEBPACK_IMPORTED_MODULE_5__app_component__["a" /* AppComponent */]]
    })
], AppModule);

//# sourceMappingURL=app.module.js.map

/***/ }),

/***/ "../../../../../src/app/app.routing.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AppRoutes; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__components_home_home_comonent__ = __webpack_require__("../../../../../src/app/components/home/home.comonent.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__components_tutorial_tutorial_comonent__ = __webpack_require__("../../../../../src/app/components/tutorial/tutorial.comonent.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__components_examples_examples_comonent__ = __webpack_require__("../../../../../src/app/components/examples/examples.comonent.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__components_about_about_comonent__ = __webpack_require__("../../../../../src/app/components/about/about.comonent.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__components_faq_faq_comonent__ = __webpack_require__("../../../../../src/app/components/faq/faq.comonent.ts");





var AppRoutes = [
    {
        path: '',
        component: __WEBPACK_IMPORTED_MODULE_0__components_home_home_comonent__["a" /* HomeComponent */],
        data: { title: 'CX Programming Language' }
    },
    {
        path: 'tutorial',
        component: __WEBPACK_IMPORTED_MODULE_1__components_tutorial_tutorial_comonent__["a" /* TutorialComponent */],
        data: { title: 'Tutorial' }
    },
    {
        path: 'examples',
        component: __WEBPACK_IMPORTED_MODULE_2__components_examples_examples_comonent__["a" /* ExamplesComponent */],
        data: { title: 'Examples' }
    },
    {
        path: 'about',
        component: __WEBPACK_IMPORTED_MODULE_3__components_about_about_comonent__["a" /* AboutComponent */],
        data: { title: 'About' }
    },
    {
        path: 'faq',
        component: __WEBPACK_IMPORTED_MODULE_4__components_faq_faq_comonent__["a" /* FAQComponent */],
        data: { title: 'FAQ' }
    },
];
//# sourceMappingURL=app.routing.js.map

/***/ }),

/***/ "../../../../../src/app/components/about/about.comonent.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AboutComponent; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__ = __webpack_require__("../../../platform-browser/@angular/platform-browser.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var AboutComponent = (function () {
    function AboutComponent(titleService, router) {
        this.titleService = titleService;
        this.router = router;
    }
    AboutComponent.prototype.ngOnInit = function () {
        this.titleService.setTitle('About');
    };
    return AboutComponent;
}());
AboutComponent = __decorate([
    Object(__WEBPACK_IMPORTED_MODULE_0__angular_core__["o" /* Component */])({
        selector: 'app-about',
        template: __webpack_require__("../../../../../src/app/components/about/about.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/about/about.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["b" /* Title */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["b" /* Title */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */]) === "function" && _b || Object])
], AboutComponent);

var _a, _b;
//# sourceMappingURL=about.comonent.js.map

/***/ }),

/***/ "../../../../../src/app/components/about/about.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".wrap {\n  margin-top: 200px;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/about/about.component.html":
/***/ (function(module, exports) {

module.exports = "\n<div class=\"page-header\" id=\"banner\">\n<div class=\"wrap\">\n    ABOUT CX.\n  </div>\n\n</div>\n"

/***/ }),

/***/ "../../../../../src/app/components/examples/examples.comonent.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return ExamplesComponent; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__ = __webpack_require__("../../../platform-browser/@angular/platform-browser.es5.js");
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var ExamplesComponent = (function () {
    function ExamplesComponent(titleService) {
        this.titleService = titleService;
    }
    ExamplesComponent.prototype.ngOnInit = function () {
        this.titleService.setTitle('Examples');
    };
    return ExamplesComponent;
}());
ExamplesComponent = __decorate([
    Object(__WEBPACK_IMPORTED_MODULE_0__angular_core__["o" /* Component */])({
        selector: 'app-examples',
        template: __webpack_require__("../../../../../src/app/components/examples/examples.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/examples/examples.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["b" /* Title */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["b" /* Title */]) === "function" && _a || Object])
], ExamplesComponent);

var _a;
//# sourceMappingURL=examples.comonent.js.map

/***/ }),

/***/ "../../../../../src/app/components/examples/examples.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/examples/examples.component.html":
/***/ (function(module, exports) {

module.exports = "\n\n\n<div class=\"page-header\" id=\"banner\">\n  <div class=\"row\">\n    <div class=\"col-lg-6 col-md-6 col-sm-6\">\n      <h1>CX Examples</h1>\n      <p class=\"lead\">Contents</p>\n\n      <ul>\n        <li><a href=\"examples#hello-world\">Hello world</a></li>\n        <li><a href=\"examples#comments\">Comments</a></li>\n        <li><a href=\"examples#definitions\">Definitions</a></li>\n        <li><a href=\"examples#functions\">Functions</a></li>\n        <li><a href=\"examples#structs\">Structs</a></li>\n        <li><a href=\"examples#arrays\">Arrays</a></li>\n        <li><a href=\"examples#expressions\">Expressions</a></li>\n        <li><a href=\"examples#go-to\">Go-to</a></li>\n        <li><a href=\"examples#if-and-if-else\">If and if/else</a></li>\n        <li><a href=\"examples#looping\">Looping</a></li>\n        <li><a href=\"examples#evolving-a-function\">Evolving a function</a></li>\n        <li><a href=\"examples#meta-programming-commands\">Meta-programming commands</a></li>\n        <li><a href=\"examples#meta-programming-functions\">Meta-programming functions</a></li>\n        <li><a href=\"examples#factorial\">Factorial</a></li>\n        <li><a href=\"examples#robot-simulator\">Robot simulator</a></li>\n      </ul>\n\n    </div>\n\n  </div>\n</div>\n\n<!-- Hello world start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"hello-world\">Hello world</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nfunc main () (out str) {{\"{\"}}\n     printStr(\"Hello World!\")\n{{\"}\"}}</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>A classic. <i>printStr</i> is part of the core module or package, and it is automatically added to any CX program, so we don't need to explicitly include it.</p>\n          <p>Something interesting in this example (as can be noted in the <a href=\"/App/Index\">CX Playground</a>) is that this program will print \"Hello World!\" two times. What's actually happening is that CX is printing to the standard output the string, and then the program is exiting, returning a string.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Hello world end -->\n\n\n\n<!-- Comments start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"comments\">Comments</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nfunc main () (out i32) {{\"{\"}}\n\t// This won't be evaluated\n\t//out := divI32(5, 0) // This won't either\n\n\tout := divI32(10, 5) // The program will return 2\n\n\t/*\n        Comment block\n\tout := subF32(3.33, 1.11)\n\tout := subF32(3.33, 1.11)\n\tout := subF32(3.33, 1.11)\n\tout := subF32(3.33, 1.11)\n        */\n{{\"}\"}}</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>Comments in CX follow the same syntax as in many other popular languages, like C and Go. Use // to comment out one line, and /* ... */ to comment out multiple lines in a source file.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Comments end -->\n\n\n\n\n<!-- Definitions start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"definitions\">Definitions</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nvar global1 i32 = 5\nvar global2 f32 = 3.14159\n\nfunc printGlobals () () {{\"{\"}}\n\tprintI32(global1)\n\tprintF32(global2)\n{{\"}\"}}\n\nfunc main () (out i32) {{\"{\"}}\n\tprintGlobals()\n\tprintF32(global2)\n\tprintI32(global1)\n{{\"}\"}}</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>A definition is just a global variable. This means that any function in the current module or in those modules that have imported the current module can access the definition's value.</p>\n          <p>In this example, we can see that both <i>printGlobals</i> and <i>main</i> are accessing the definitions <i>global1</i> and <i>global2</i>.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Definitions end -->\n\n\n\n\n\n<!-- Functions start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"functions\">Functions</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nvar PI f32 = 3.14159\n\nfunc circleArea (radius f32) (area f32) {{\"{\"}}\n\tmulF32(mulF32(radius, radius), PI)\n{{\"}\"}}\n\nfunc circlePerimeter (radius f32) (perimeter f32) {{\"{\"}}\n\tperimeter := mulF32(mulF32(2, radius), PI)\n{{\"}\"}}\n\nfunc main () (area f32) {{\"{\"}}\n\tarea := circleArea(2)\n\tprintF32(circlePerimeter(5))\n{{\"}\"}}</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>The example defines two functions: one for calculating a circle's area, and another for calculating a circle's perimeter. A feature in CX is that the output of a function can be explicitly defined, as in the case of <i>circlePerimeter</i>. If the output is not explicit, CX will take the last expression as the output of the function.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Functions end -->\n\n\n<!-- Structs start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"structs\">Structs</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\ntype Point struct {{\"{\"}}\n\tname str\n\tx i32\n\ty i32\n{{\"}\"}}\n\nfunc main () (out i32) {{\"{\"}}\n\tvar uPoint Point\n\tprintStr(uPoint.name)\n\tprintI32(uPoint.x)\n\tprintI32(uPoint.y)\n\n\tvar myPoint Point\n\tmyPoint.x = addI32(5, 10)\n\tmyPoint.y = addI32(myPoint.x, 20)\n\n\tout = addI32(myPoint.x, myPoint.y)\n{{\"}\"}}\n            </code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>User-defined types can be created by the use of structs. In this example, a type called \"Point\" is created. We want a point to be defined by a name, and its x and y coordinates on a plane.</p>\n          <p>The <i>main</i> function defines two variables of type Point, one at <b>line 12</b>, and another at <b>line 17</b>. <b>Lines 11, 12, and 13</b> try to print the values of an uninitialized instance of Point, but this won't crash the program, as CX initializes the fields to their corresponding zero value. For example, a <i>str</i> will be initialized to an empty string, a <i>bool</i> will be initialized to false, and an <i>i32</i> to a 32 bit integer 0.</p>\n          <p>At <b>line 15</b> another instance of Point is created, and their x and y fields are initialized at <b>lines 16 and 16</b>. Finally, at <b>line 19</b>, we perform the addition of both fields and its sum serves as the output of the program.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Structs end -->\n\n\n<!-- Arrays start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"arrays\">Arrays</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nfunc main () (out str) {{\"{\"}}\n\tvar ids []i32 = []i32{{\"{\"}}433, 561, 652, 984{{\"}\"}}\n\tvar ages []i32 = []i32{{\"{\"}}23, 21, 26, 31{{\"}\"}}\n\tvar grades []f32 = []f32{{\"{\"}}8.8, 9.4, 9.3, 8.3{{\"}\"}}\n\n\tprintStr(\"Student's ID:\")\n\tprintI32(readI32A(ids, 1))\n\tprintStr(\"Student's age:\")\n\tprintI32(readI32A(ages, 1))\n\tprintStr(\"Student's age:\")\n\tprintF32(readF32A(grades, 1))\n\tprintStr(\"done.\")\n{{\"}\"}}</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>Three arrays are created to simulate the data of four students. They hold the ids, ages and grades. At <b>lines 9, 11, and 13</b> array reader functions are used to access the values at index 1 for each of the arrays.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Arrays end -->\n\n\n\n<!-- Expressions start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"expressions\">Expressions</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nfunc sayHi () () {{\"{\"}}\n\tprintStr(\"Hi\")\n{{\"}\"}}\n\nfunc sayMyName (name str) () {{\"{\"}}\n\tprintStr(name)\n{{\"}\"}}\n\nfunc returnName () (name str) {{\"{\"}}\n\tname := idStr(\"Bart\")\n{{\"}\"}}\n\nfunc multiReturn (num1 i32, num2 i32) (add i32, sub i32, mul i32, div i32) {{\"{\"}}\n\tadd := addI32(num1, num2)\n\tsub := subI32(num1, num2)\n\tmul := mulI32(num1, num2)\n\tdiv := divI32(num1, num2)\n{{\"}\"}}\n\nfunc main () (out i32) {{\"{\"}}\n\tsayHi()\n\tsayMyName(\"Homer\")\n\tprintStr(returnName())\n\ta, s, out, d := multiReturn(20, 20)\n{{\"}\"}}</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>Four functions and a main function are defined in this example. <i>sayHi</i> demonstrates the definition of an input-less and output-less function. <i>sayMyName</i> receives one parameter as input, but does not return anything. <i>returnName</i> receives no parameters as input, but returns one output.</p>\n\n          <p><i>multiReturn</i> is a bit more interesting than the other functions, as it receives two input parameters, and returns four output parameters. This function takes two 32 bit integers and returns their addition, subtraction, multiplication and division. The main function calls each of these functions, and assigns its output with the third output parameter from <i>multiReturn</i>.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Expressions end -->\n\n\n<!-- Go-to start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"go-to\">Go-to</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nfunc basicIf (num i32) (num i32) {{\"{\"}}\n\tpred := gtI32(num, 0)\n\tgoTo(pred, 1, 3)\n\tprintStr(\"Greater than 0\")\n\tgoTo(true, 10, 0)\n\tprintStr(\"Less than 0\")\n{{\"}\"}}\n\nfunc main () (out i32) {{\"{\"}}\n\tbasicIf(5)\n{{\"}\"}}</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p><i>goTo</i> <b>should not</b> be used, as it can create programs that are very hard to debug. CX uses the go-to operator to create other flow-control structures, such as <i>if</i> and <i>for</i>.</p>\n          <p>Nevertheless, if you really require a go-to in your program, you can use it. The example above shows how an if-else can be constructed using two go-to expressions. <i>goTo</i> takes three arguments as its parameters: a predicate, a number of lines the program should advance if the predicate evaluates to true, and a number of lines the program should advance if the predicate evaluates to false.</p>\n          <p>As in the case of the second <i>goTo</i> in the example, if you give it a number of lines which exceeds the number of expressions in the function, this will only cause the function to return. Negative numbers can also be provided, and they will cause the program to re-evaluate the N expressions above.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Go-to end -->\n\n\n<!-- If-and-if-else start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"if-and-if-else\">If and if/else</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nfunc main () (out i32) {{\"{\"}}\n     \tif false {{\"{\"}}\n\t\terror := divI32(50, 0)\n\t\tprintStr(\"This will never be printed\")\n\t{{\"}\"}}\n\n\tif true {{\"{\"}}\n\t\tprintStr(\"This will always print\")\n\t{{\"}\"}}\n\n\tif gtI32(5, 3) {{\"{\"}}\n\t\tprintStr(\"5 is greater than 3\")\n\t{{\"}\"}}\n\n\tif eqStr(\"password123\", \"password123\") {{\"{\"}}\n\t\tprintStr(\"Access granted\")\n\t{{\"}\"}} else {{\"{\"}}\n\t\tprintStr(\"Access denied\")\n\t{{\"}\"}}\n\n\tif lteqI32(50, 5) {{\"{\"}}\n\t\tout = idI32(100)\n\t{{\"}\"}} else {{\"{\"}}\n\t\tout = idI32(200)\n\t{{\"}\"}}\n{{\"}\"}}</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>If and if/else statements in CX work as in any other programming language. If the condition evaluates to true, the first block of expressions will be executed; if the condition evaluates to false, the second block of expressions will be executed.</p>\n          <p>Like in Go, the condition doesn't need to be enclosed in parenthesis.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- If-and-if-else end -->\n\n\n\n<!-- Looping start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"looping\">Looping</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nfunc main () (out i32) {{\"{\"}}\n\tvar out i32 = 0\n\tfor ltI32(out, 10) {{\"{\"}}\n\t\tprintI32(out)\n\t\tout = addI32(out, 1)\n\t{{\"}\"}}\n{{\"}\"}}</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>You can use <i>for</i> to create a loop. As in other programming languages that implement a for statement, a condition is provided which will instruct the computer to evaluate a block of expressions while the condition evaluates to true. In the case of the example above, the for loop will print the numbers from 0 to 10.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Looping end -->\n\n\n\n<!-- Evolving a function start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"evolving-a-function\">Evolving a function</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nvar inps []f64 = []f64{{\"{\"}}\n\t-10.0, -9.0, -8.0, -7.0, -6.0, -5.0, -4.0, -3.0, -2.0, -1.0,\n\t0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0{{\"}\"}}\n\nvar outs []f64 = []f64{{\"{\"}}\n\t90.0, 72.0, 56.0, 42.0, 30.0, 20.0, 12.0, 6.0, 2.0, 0.0, 0.0,\n\t2.0, 6.0, 12.0, 20.0, 30.0, 42.0, 56.0, 72.0, 90.0, 110.0{{\"}\"}}\n\nfunc solution (n f64) (out f64) {{\"{\"}}\n\tsquare = mulF64(n, n)\n{{\"}\"}}\n\nfunc main () (out f64) {{\"{\"}}\n\t:dStack false;\n\t:dProgram true;\n\tevolve(\"solution\", \"addF64|mulF64|subF64\", inps, outs, 5, 300, f32ToF64(0.1))\n\n\tprintStr(\"Extrapolating solution\")\n\tprintF64(solution(f32ToF64(30.0)))\n{{\"}\"}}</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p><i>evolve</i> is perhaps the most complicated native function in CX, not only because it is the native function that receives the greater number of parameters, but because of its purpose.</p>\n          <p><i>evolve</i> follows the principles of evolutionary computation. In particular, <i>evolve</i> performs a technique called genetic programming. Genetic programming tries to find a combination of operators and arguments that will solve a problem. For example, you could instruct evolve to find a combination of operators that, when sent 10, returns 20. This might sound trivial, but genetic programming and other evolutionary algorithms can solve very complicated problems.</p>\n          <p>In this example, we have a number of inputs defined at <b>line 3</b>, and we want to find a function which maps them to the outputs defined at <b>line 7</b>. The function <i>solution</i> could be totally empty, but we can help <i>evolve</i> a little bit. We know that a close solution to the problem is to multiply n by n, so we add that expression.</p>\n          <p>The call to <i>evolve</i> at <b>line 18</b> takes \"solution\" as its first argument. This instructs the evolutionary algorithm to evolve the <i>solution</i> function. Next, we provide the function with a set of functions to choose from when the algorithm starts creating different combinations. This argument is actually a regular expression, so we could send it \".\" to match any function in the current program. However, in genetic programming it's always useful to limit the algorithm to a small set of operators. If we have a very large set of operators, the algorithm could take too long to find a solution.</p>\n          <p>The next two parameters are the set of inputs and outputs we want to create a function for. Then, the next parameter defines the number of expressions we want the target solution function to be limited to. In traditional genetic programming a limit doesn't exist, and the resulting programs are usually very large. CX's <i>evolve</i> follows a strategy used in cartesian genetic programming (CGP). In CGP, a limit of statements or expressions is given, and this eliminates bloat in the solutions.</p>\n          <p>The last two parameters are know as the stop criteria. Evolutionary algorithms are heuristics, which means that they won't necessarily deliver the optimal solution to a problem. By using the stop criteria, we can tell <i>evolve</i> to give us the best solution it could find after N iterations, or if the error obtained is less than certain threshold.</p>\n          <p>If you run the example above, you should get a list of errors printed on the screen. These errors will decrease as the evolution algorithm progresses. When either the error threshold or the number of iterations is reached, the evolutionary algorithm will stop and the program will print what the evolved solution function evaluates to when sent an argument of 30.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Evolving a function end -->\n\n\n\n<!-- Meta-programming commands start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"meta-programming-commands\">Meta-programming commands</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">* :dStack false;\n* :dProgram true;\n\nProgram\n0.- Module: main\n\tDefinitions\n\t\t0.- Definition: inps []f64\n\t\t1.- Definition: outs []f64\n\tFunctions\n\t\t0.- Function: solution (n f64) (out f64)\n\t\t\t0.- Expression: square = mulF64(n f64, n f64)\n\t\t\t1.- Expression: var_8 = addF64(n f64, square )\n\t\t\t2.- Expression: var_300 = addF64(square , square )\n\t\t\t3.- Expression: var_303 = mulF64(square , var_8 )\n\t\t\t4.- Expression: out = addF64(n f64, square )\n\t\t1.- Function: main () (out f64)\n\t\t\t0.- Expression: nonAssign_0 = f32ToF64(0.1 f32)\n\t\t\t1.- Expression: nonAssign_2 = printStr(\"Extrapolating solution\" str)\n\t\t\t2.- Expression: nonAssign_3 = f32ToF64(30 f32)\n\t\t\t3.- Expression: nonAssign_4 = solution(nonAssign_3 )\n\t\t\t4.- Expression: nonAssign_5 = printF64(nonAssign_4 )\n\n* :step 100;\nEvolving function 'solution'\n5.238095238095238\n5.238095238095238\n5.238095238095238\n5.238095238095238\n5.238095238095238\n5.238095238095238\n5.238095238095238\n5.238095238095238\n5.238095238095238\n5.238095238095238\n0\nFinished evolving function 'solution'\nExtrapolating solution\n930</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>This example is not meant to be run as a program. It's actually a simulation of a REPL session. The reason behind this is that meta-programming commands (unlike meta-programming functions) are most useful when used in an interactive environment.</p>\n          <p>In this example, three meta-programming commands are used: dStack, dProgram, and step. The program that was loaded into the REPL is the example from <a href=\"evolving-a-function\">evolving a function</a>. <i>:dStack</i> is used for telling the REPL that it should start/stop printing the call stack each time the program advances. :dProgram is like :dStack, but it tells the REPL to print the program's abstract syntax tree. The \"d\" at the start of these commands stands for \"debug\", and this tells us that both of these commands are used for debugging purposes only.</p>\n          <p>:step takes an integer as its argument, and tells CX to run forwards or backwards the number of steps indicated by this integer. It is worth noting that giving :step a negative number is not the same as, for example, using <i>goTo</i> to get back to a previous expression in a function. :step will \"forget\" anything that happened in the last N steps, if N is negative. If N is positive, it's the equivalent of a normal execution of a CX program.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Meta-programming commands end -->\n\n<!-- Meta-programming functions start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"meta-programming-functions\">Meta-programming functions</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nvar greeting str = \"Meta hello.\"\nvar farewell str = \"Meta good-bye.\"\n\nfunc main () () {{\"{\"}}\n\tremExpr(\"toBeRemoved\")\n\t:tag toBeRemoved;\n\tdivI32(3, 0)\n\taddExpr(\"newPrint\", \"printStr\")\n\taffExpr(\"newPrint\", \".\", 0)\n{{\"}\"}}</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>In this example we can see one of the special features of CX: meta-programming functions. Unlike meta-programming commands, which unleash their true potential only in the REPL, meta-programming functions can be used in CX expressions and modify a program's structure in runtime.</p>\n          <p>We can see that at <b>line 8</b>, a division by 0 is defined, but the program won't crash because we are removing this expression by calling <i>remExpr(\"main\", 1)</i> at <b>line 7</b>, i.e., we are telling CX to remove expression 1 (the second one, as expressions are zero-indexed) from the \"main\" function.</p>\n          <p>At <b>line 9</b>, we instruct CX to add an <i>printStr</i> expression. This call to <i>addExpr</i> would normally crash the program, as it's adding an argument-less call to the function. However, we then ask CX to add whatever it wants as an argument to the expression we just added, at <b>line 10</b>. The \".\" (a regular expression) we are sending to <i>exprAff</i> is telling CX to apply any affordance it wants.</p>\n          <p>As a more elaborated example of affordances and meta-programming functions, check the <a href=\"robot-simulator\">robot simulator</a>.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Meta-programming functions end -->\n\n<!-- Factorial start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"factorial\">Factorial</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nfunc factorial (num i32) (fact i32) {{\"{\"}}\n\tif eqI32(num, 1) {{\"{\"}}\n\t\tfact := idI32(1)\n\t{{\"}\"}} else {{\"{\"}}\n\t\tfact := mulI32(num, factorial(subI32(num, 1)))\n\t{{\"}\"}}\n{{\"}\"}}\n\nfunc main () (out i32) {{\"{\"}}\n\tfactorial(6)\n{{\"}\"}}</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>The <i>factorial</i> function defined above shows how one can use recursion in CX. If we give the compiler/interpreter the :dStack meta-programming command, we'd see how CX calls <i>factorial</i> several times, and how each of these calls are waiting for the next to return.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Factorial end -->\n\n\n\n\n\n\n<!-- Robot simulator start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"robot-simulator\">Robot simulator</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n            <pre><code class=\"go\">package main\n\nvar goNorth str = \"going north.\"\nvar goSouth str = \"going south.\"\nvar goWest str = \"going west.\"\nvar goEast str = \"going east.\"\n\nfunc robotAct (row i32, col i32, action str) (row i32, col i32) {{\"{\"}}\n\tprintStr(action)\n\tif eqStr(action, \"going north.\") {{\"{\"}}\n\t\trow = addI32(row, -1)\n\t{{\"}\"}}\n\tif eqStr(action, \"going south.\") {{\"{\"}}\n\t\trow = addI32(row, 1)\n\t{{\"}\"}}\n\tif eqStr(action, \"going east.\") {{\"{\"}}\n\t\tcol = addI32(col, 1)\n\t{{\"}\"}}\n\tif eqStr(action, \"going west.\") {{\"{\"}}\n\t\tcol = addI32(col, -1)\n\t{{\"}\"}}\n{{\"}\"}}\n\nfunc robot (row i32, col i32) (row i32, col i32) {{\"{\"}}\n\tremArg(\"robotAct\")\n\taffExpr(\"robotAct\", \"goNorth|goSouth|goWest|goEast\", 0)\n\t:tag robotAct;\n\trow, col = robotAct(row, col, \"\")\n{{\"}\"}}\n\nfunc map2Dto1D (r i32, c i32, w i32) (i i32) {{\"{\"}}\n\ti = addI32(mulI32(w, r), c)\n{{\"}\"}}\n\nfunc map1Dto2D (i i32, w i32) (r i32, c i32) {{\"{\"}}\n\tr = divI32(i, W)\n\tc = modI32(i, W)\n{{\"}\"}}\n\nfunc robotObjects (row i32, col i32, width i32, wallMap []bool, wormholeMap []bool) () {{\"{\"}}\n\tremObjects()\n\tif readBoolA(wallMap, map2Dto1D(addI32(row, -1), col, width)) {{\"{\"}}\n\t\taddObject(\"northWall\")\n\t{{\"}\"}}\n\tif readBoolA(wallMap, map2Dto1D(addI32(row, 1), col, width)) {{\"{\"}}\n\t\taddObject(\"southWall\")\n\t{{\"}\"}}\n\tif readBoolA(wallMap, map2Dto1D(row, addI32(col, 1), width)) {{\"{\"}}\n\t\taddObject(\"eastWall\")\n\t{{\"}\"}}\n\tif readBoolA(wallMap, map2Dto1D(row, addI32(col, -1), width)) {{\"{\"}}\n\t\taddObject(\"westWall\")\n\t{{\"}\"}}\n\n\tif readBoolA(wormholeMap, map2Dto1D(addI32(row, -1), col, width)) {{\"{\"}}\n\t\taddObject(\"northWormhole\")\n\t{{\"}\"}}\n\tif readBoolA(wormholeMap, map2Dto1D(addI32(row, 1), col, width)) {{\"{\"}}\n\t\taddObject(\"southWormhole\")\n\t{{\"}\"}}\n\tif readBoolA(wormholeMap, map2Dto1D(row, addI32(col, 1), width)) {{\"{\"}}\n\t\taddObject(\"eastWormhole\")\n\t{{\"}\"}}\n\tif readBoolA(wormholeMap, map2Dto1D(row, addI32(col, -1), width)) {{\"{\"}}\n\t\taddObject(\"westWormhole\")\n\t{{\"}\"}}\n{{\"}\"}}\n\nfunc main () (out str) {{\"{\"}}\n\n\tsetClauses(\"\n          aff(robotAct, goNorth, X, R) :- X = northWall, R = false.\n\t  aff(robotAct, goSouth, X, R) :- X = southWall, R = false.\n\t  aff(robotAct, goWest, X, R) :- X = westWall, R = false.\n\t  aff(robotAct, goEast, X, R) :- X = eastWall, R = false.\n\n\t  aff(robotAct, goNorth, X, R) :- X = northWormhole, R = true.\n\t  aff(robotAct, goSouth, X, R) :- X = southWormhole, R = true.\n\t  aff(robotAct, goWest, X, R) :- X = westWormhole, R = true.\n\t  aff(robotAct, goEast, X, R) :- X = eastWormhole, R = true.\n        \")\n\n\tsetQuery(\"aff(%s, %s, %s, R).\")\n\n\tvar wallMap []bool = []bool{{\"{\"}}\n\t\ttrue, true,  true,  true,  true,\n\t\ttrue, false, true, false, true,\n\t\ttrue, false, true, false, true,\n\t\ttrue, false, false, false, true,\n\t\ttrue, true,  true,  true,  true{{\"}\"}}\n\n\tvar wormholeMap []bool = []bool{{\"{\"}}\n\t\tfalse, false, false, false, false,\n\t\tfalse, false, false, false, false,\n\t\tfalse, false, false, false, false,\n\t\tfalse, false, false, false, false,\n\t\tfalse, false, false, false, false{{\"}\"}}\n\n\tvar width i32 = 5\n\tvar row i32 = 1\n\tvar col i32 = 1\n\n\tvar counter i32\n\tfor ltI32(counter, 6) {{\"{\"}}\n\t\twallMap = writeBoolA(wallMap, map2Dto1D(row, col, width), true)\n\t\twormholeMap = writeBoolA(wormholeMap, map2Dto1D(row, col, width), false)\n\t\trobotObjects(row, col, width, wallMap, wormholeMap)\n\t\trow, col := robot(row, col)\n\t\tcounter = addI32(counter, 1)\n\t{{\"}\"}}\n\n\tprintStr(\"done.\")\n{{\"}\"}}\n</code></pre>\n        </div>\n        <div class=\"panel-body\">\n          <p>This example tries to illustrate most of the concepts shown in the previous examples.</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Robot simulator end -->\n\n\n"

/***/ }),

/***/ "../../../../../src/app/components/faq/faq.comonent.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return FAQComponent; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__ = __webpack_require__("../../../platform-browser/@angular/platform-browser.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var FAQComponent = (function () {
    function FAQComponent(titleService, router) {
        this.titleService = titleService;
        this.router = router;
    }
    FAQComponent.prototype.ngOnInit = function () {
        this.titleService.setTitle('FAQ');
    };
    return FAQComponent;
}());
FAQComponent = __decorate([
    Object(__WEBPACK_IMPORTED_MODULE_0__angular_core__["o" /* Component */])({
        selector: 'app-faq',
        template: __webpack_require__("../../../../../src/app/components/faq/faq.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/faq/faq.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["b" /* Title */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["b" /* Title */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */]) === "function" && _b || Object])
], FAQComponent);

var _a, _b;
//# sourceMappingURL=faq.comonent.js.map

/***/ }),

/***/ "../../../../../src/app/components/faq/faq.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".wrap {\n  margin-top: 200px;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/faq/faq.component.html":
/***/ (function(module, exports) {

module.exports = "\n<div class=\"page-header\" id=\"banner\">\n<div class=\"wrap\">\n    CX FAQ.\n  </div>\n</div>\n"

/***/ }),

/***/ "../../../../../src/app/components/header/header.comonent.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return HeaderComponent; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};

var HeaderComponent = (function () {
    function HeaderComponent() {
        this.isCollapsed = true;
    }
    HeaderComponent.prototype.collapsed = function (event) {
        console.log(event);
    };
    HeaderComponent.prototype.expanded = function (event) {
        console.log(event);
    };
    return HeaderComponent;
}());
HeaderComponent = __decorate([
    Object(__WEBPACK_IMPORTED_MODULE_0__angular_core__["o" /* Component */])({
        selector: 'app-header',
        template: __webpack_require__("../../../../../src/app/components/header/header.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/header/header.component.css")]
    })
], HeaderComponent);

//# sourceMappingURL=header.comonent.js.map

/***/ }),

/***/ "../../../../../src/app/components/header/header.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".hover:hover {\n  cursor: pointer;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/header/header.component.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"navbar navbar-inverse navbar-fixed-top\">\n  <div class=\"container\">\n    <div class=\"navbar-header\">\n      <a [routerLink]=\"['']\" class=\"navbar-brand\">CX Programming Language </a>\n      <button class=\"navbar-toggle\" type=\"button\" (click)=\"isCollapsed = !isCollapsed\">\n        <span class=\"icon-bar\"></span>\n        <span class=\"icon-bar\"></span>\n        <span class=\"icon-bar\"></span>\n      </button>\n    </div>\n     <div class=\"navbar-collapse collapse\"\n         (collapsed)=\"collapsed($event)\"\n         (expanded)=\"expanded($event)\"\n         [collapse]=\"isCollapsed\">\n      <ul class=\"nav navbar-nav\">\n        <li class=\"dropdown\" dropdown>\n          <a class=\"dropdown-toggle hover\" dropdownToggle >About <span class=\"caret\"></span></a>\n          <ul class=\"dropdown-menu\" *dropdownMenu  aria-labelledby=\"themes\">\n            <li><a routerLink=\"about\">About CX</a></li>\n            <li class=\"divider\"></li>\n            <li><a routerLink=\"faq\">FAQ</a></li>\n            <li><a href=\"https://www.skycoin.net/\">Skycoin</a></li>\n          </ul>\n        </li>\n        <li class=\"dropdown\" dropdown>\n          <a class=\"dropdown-toggle hover\" dropdownToggle  aria-expanded=\"false\" >CX <span class=\"caret\"></span></a>\n          <ul class=\"dropdown-menu\" *dropdownMenu aria-labelledby=\"download\">\n            <li><a [routerLink]=\"['']\">CX Playground</a></li>\n            <li class=\"divider\"></li>\n            <li><a routerLink=\"tutorial\">Tutorial and Documentation</a></li>\n            <li><a routerLink=\"examples\">Examples</a></li>\n            <li><a href=\"/assets/cx.pdf\">Specification</a></li>\n          </ul>\n        </li>\n        <li>\n          <a href=\"https://blog.skycoin.net/\">Blog</a>\n        </li>\n      </ul>\n\n      <ul class=\"nav navbar-nav navbar-right\">\n        <li><a [routerLink]=\"['']\">CX Playground</a></li>\n        <li><a href=\"https://github.com/skycoin/cx\" target=\"_blank\">Github Repository</a></li>\n      </ul>\n\n    </div>\n  </div>\n</div>\n\n\n\n"

/***/ }),

/***/ "../../../../../src/app/components/home/home.comonent.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return HomeComponent; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__ = __webpack_require__("../../../platform-browser/@angular/platform-browser.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__services_api_service__ = __webpack_require__("../../../../../src/app/services/api.service.ts");
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};




var HomeComponent = (function () {
    function HomeComponent(titleService, router, api) {
        this.titleService = titleService;
        this.router = router;
        this.api = api;
        this.programms = [
            { id: 1, name: 'Hello world', code: 'package main \n\n func main () (){\n \tstr.print("Hello World!")\n}' },
            { id: 2, name: 'Looping', code: 'package main\r\n\r\nfunc main () () {\r\n\tfor c := 0; i32.lt(c, 20); c = i32.add(c, 1) {\r\n\t\ti32.print(c)\r\n\t}\r\n}' },
            { id: 3, name: 'Factorial', code: 'package main\r\n\r\nfunc factorial (num i32) (fact i32) {\r\n\tif i32.eq(num, 1) {\r\n\t\tfact = 1\r\n\t} else {\r\n\t\tfact = i32.mul(num, factorial(i32.sub(num, 1)))\r\n\t}\r\n}\r\n\r\nfunc main () () {\r\n\ti32.print(factorial(6))\r\n}' },
            { id: 4, name: 'Evolving a function', code: 'package main\r\n\r\nvar inps []f64 = []f64{\r\n\t-10.0, -9.0, -8.0, -7.0, -6.0, -5.0, -4.0, -3.0, -2.0, -1.0,\r\n\t0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}\r\n\r\nvar outs []f64 = []f64{\r\n\t90.0, 72.0, 56.0, 42.0, 30.0, 20.0, 12.0, 6.0, 2.0, 0.0, 0.0,\r\n\t2.0, 6.0, 12.0, 20.0, 30.0, 42.0, 56.0, 72.0, 90.0, 110.0}\r\n\r\nfunc solution (n f64) (out f64) {}\r\n\r\nfunc main () (out f64) {\r\n\tevolve(\"solution\", \"f64.add|f64.mul|f64.sub\", inps, outs, 2, 100, f32.f64(0.1))\r\n\r\n\tstr.print(\"Extrapolating evolved solution\")\r\n\tf64.print(solution(f32.f64(30.0)))\r\n}' },
            { id: 5, name: 'Text-based adventure', code: "package main\r\n\r\nfunc walk (flag bool) () {\r\n\tif flag {\r\n\t\tif i32.gt(i32.rand(0, 10), 5) {\r\n\t\t\tstr.print(\"The traveler keeps following the lane, making sure to ignore any pain.\")\r\n\t\t} else {\r\n\t\t\tstr.print(\"North, east, west, south. Any direction is good, as long as no monster can be found.\")\r\n\t\t}\r\n\t}\r\n}\r\n\r\nfunc noise (flag bool) () {\r\n\tif flag {\r\n\t\tif i32.gt(i32.rand(0, 10), 5) {\r\n\t\t\tstr.print(\"A cracking noise is heard, but no monster is there.\")\r\n\t\t} else {\r\n\t\t\taddObject(\"monster\")\r\n\t\t\tstr.print(\"Howling and growling, the monster is coming.\")\r\n\t\t}\r\n\t}\r\n}\r\n\r\nfunc consider (flag bool) () {\r\n\tif flag {\r\n\t\tif i32.gt(i32.rand(0, 10), 5) {\r\n\t\t\tremObject(\"monster\")\r\n\t\t\tstr.print(\"The traveler runs away, and cowardice lets him live for another day.\")\r\n\t\t} else {\r\n\t\t\taddObject(\"fight\")\r\n\t\t\tstr.print(\"Bravery comes into sight, in the hope of living for another night.\")\r\n\t\t}\r\n\t}\r\n}\r\n\r\nfunc chance (flag bool) () {\r\n\tif and(flag, i32.gt(i32.rand(0, 10), 5)) {\r\n\t\tif i32.gt(i32.rand(0, 10), 5) {\r\n\t\t\tremObject(\"fight\")\r\n\t\t\tremObject(\"monster\")\r\n\t\t\tstr.print(\"The monster stares, almost as in compassion, and leaves despite the traveler's past actions.\")\r\n\t\t} else {\r\n\t\t\tremObject(\"fight\")\r\n\t\t\tstr.print(\"The monster starts a deep glare, waiting for the traveler to accept the dare.\")\r\n\t\t}\r\n\t}\r\n}\r\n\r\nfunc fightResult (flag bool) () {\r\n\tif flag {\r\n\t\tif i32.gt(i32.rand(0, 10), 5) {\r\n\t\t\taddObject(\"died\")\r\n\t\t\tstr.print(\"But failure describes this fend and, suddenly, this adventure comes to an end.\")\r\n\t\t} else {\r\n\t\t\tremObject(\"monster\")\r\n\t\t\tremObject(\"fight\")\r\n\t\t\tstr.print(\"Naive, and even dumb, but the traveler's act leaves the monster numb.\")\r\n\t\t}\r\n\t}\r\n}\r\n\r\nfunc theEnd (flag bool) () {\r\n\tif flag {\r\n\t\tstr.print(\"\")\r\n\t\thalt(\"You died.\")\r\n\t}\r\n}\r\n\r\nfunc act () () {\r\n\tyes := true\r\n\tno := false\r\n\t\r\n\tremArg(\"walk\")\r\n\taffExpr(\"walk\", \"yes|no\", 0)\r\n\twalk:\r\n\twalk(false)\r\n\r\n\tremArg(\"noise\")\r\n\taffExpr(\"noise\", \"yes|no\", 0)\r\n\tnoise:\r\n\tnoise(false)\r\n\r\n\tremArg(\"consider\")\r\n\taffExpr(\"consider\", \"yes|no\", 0)\r\n\tconsider:\r\n\tconsider(false)\r\n\r\n\tremArg(\"chance\")\r\n\taffExpr(\"chance\", \"yes|no\", 0)\r\n\tchance:\r\n\tchance(false)\r\n\r\n\tremArg(\"fightResult\")\r\n\taffExpr(\"fightResult\", \"yes|no\", 0)\r\n\tfightResult:\r\n\tfightResult(false)\r\n\r\n\tremArg(\"theEnd\")\r\n\taffExpr(\"theEnd\", \"yes|no\", 0)\r\n\ttheEnd:\r\n\ttheEnd(false)\r\n}\r\n\r\nfunc main () () {\r\n\tsetClauses(\"\r\n          aff(walk, yes, X, R) :- X = monster, R = false.\r\n          aff(noise, yes, X, R) :- X = monster, R = false.\r\n\r\n          aff(consider, yes, X, R) :-  R = false.\r\n          aff(chance, yes, X, R) :-  R = false.\r\n          aff(fightResult, yes, X, R) :-  R = false.\r\n          aff(theEnd, yes, X, R) :-  R = false.\r\n\r\n          aff(consider, yes, X, R) :- X = monster, R = true.\r\n          aff(chance, yes, X, R) :- X = fight, R = true.\r\n          aff(fightResult, yes, X, R) :- X = fight, R = true.\r\n          aff(theEnd, yes, X, R) :- X = died, R = true.\r\n        \")\r\n\t\r\n\tsetQuery(\"aff(%s, %s, %s, R).\")\r\n\r\n\taddObject(\"start\")\r\n\tfor c := 0; i32.lt(c, 5); c = i32.add(c, 1) {\r\n\t\tact()\r\n\t}\r\n\r\n\tstr.print(\"\")\r\n\tstr.print(\"You survived.\")\r\n}" },
            { id: 6, name: 'More examples!', code: '' }
        ];
        this.selectedValue = this.programms[0];
        this.code = this.programms[0].code;
        this.showResult = false;
        this.result = 'waiting...';
    }
    HomeComponent.prototype.ngOnInit = function () {
        this.titleService.setTitle('CX Programming Language');
    };
    HomeComponent.prototype.changeCode = function () {
        if (this.selectedValue.id === 6) {
            this.router.navigate(['examples']);
        }
        else {
            this.code = this.selectedValue.code;
        }
    };
    HomeComponent.prototype.clearCode = function () {
        this.code = '';
    };
    HomeComponent.prototype.runCode = function () {
        var _this = this;
        console.log(this.code);
        var str = this.code;
        str = str.replace(new RegExp('\n', 'g'), ' ');
        str = str.replace(new RegExp('\t', 'g'), ' ');
        str = str.replace(new RegExp('"', 'g'), '\"');
        console.log(str);
        this.api.sendCode(str).subscribe(function (data) {
            _this.result = data._body;
            _this.showResult = true;
        });
    };
    return HomeComponent;
}());
HomeComponent = __decorate([
    Object(__WEBPACK_IMPORTED_MODULE_0__angular_core__["o" /* Component */])({
        selector: 'app-home',
        template: __webpack_require__("../../../../../src/app/components/home/home.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/home/home.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["b" /* Title */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["b" /* Title */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["a" /* Router */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_3__services_api_service__["a" /* ApiService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_3__services_api_service__["a" /* ApiService */]) === "function" && _c || Object])
], HomeComponent);

var _a, _b, _c;
//# sourceMappingURL=home.comonent.js.map

/***/ }),

/***/ "../../../../../src/app/components/home/home.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/home/home.component.html":
/***/ (function(module, exports) {

module.exports = "\n<div class=\"page-header\" id=\"banner\">\n<div class=\"row\">\n    <div class=\"col-lg-6 col-md-6 col-sm-6\">\n      <h1>CX</h1>\n      <p class=\"lead\">A strongly typed programming language</p>\n\n      <div class=\"row\">\n        <div class=\"col-lg-6 col-md-6 col-sm-6\">\n          <div class=\"list-group table-of-contents\">\n            <a class=\"list-group-item\" href=\"#compiled-and-interpreted\">Compiled and interpreted</a>\n            <a class=\"list-group-item\" href=\"#strict-typing-system\">Strict typing system</a>\n            <a class=\"list-group-item\" href=\"#affordances\">Affordances</a>\n            <a class=\"list-group-item\" href=\"#serialization\">Serialization</a>\n            <a class=\"list-group-item\" href=\"#stepping\">Stepping</a>\n            <a class=\"list-group-item\" href=\"#integrated-evolution-of-functions\">Integrated evolution of functions</a>\n            <a class=\"list-group-item\" href=\"#interactive-debugging\">Interactive debugging</a>\n          </div>\n        </div>\n      </div>\n    </div>\n\n    <div class=\"col-lg-6 col-md-6 col-sm-6\">\n      <h1>CX Playground</h1>\n      <p class=\"lead\">Try CX online</p>\n      <form action=\"/App/Index\" method=\"post\">\n        <label for=\"example-select\">Load example: </label>\n        <select id=\"example-select\" class=\"form-control\" name=\"example-select\" [(ngModel)]=\"selectedValue\"  (change)=\"changeCode()\">\n          <option *ngFor=\"let c of programms\" [selected]=\"c.name == 'Hello world'\" [ngValue]=\"c\">{{c.name}}</option>\n        </select>\n        <br />\n        <textarea id=\"code-editor\" class=\"form-control\" name=\"code\" rows=\"6\" [(ngModel)]=\"code\"></textarea>\n        <br />\n        <button type=\"submit\" class=\"btn btn-primary\" (click)=\"runCode()\">Run</button>\n        <button id=\"clear-button\" class=\"btn btn-warning\" (click)=\"clearCode()\">Clear</button>\n        <br />\n        <br />\n\n        <div class=\"alert alert-success\" *ngIf=\"showResult\">{{result}}</div>\n\n      </form>\n    </div>\n  </div>\n</div>\n\n<!-- Compiled and interpreted start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"compiled-and-interpreted\">Compiled and interpreted</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n      <p>The CX specification enforces a CX dialect to provide the developer with both an interpreter and a compiler. An interpreted program is far slower than its compiled counterpart, as is expected, but will allow a more flexible program. This flexibility comes from meta-programming functions, and affordances, which can modify a program’s structure during runtime.\n\n      </p>\n      <p>A compiled program needs a more rigid structure than an interpreted program, as many of the optimizations leverage this rigidity. As a consequence, the affordance system and any function that operates over the program’s structure will be limited in functionality in a compiled program.\n      </p>\n      <p>The compiler should be used when performance is the biggest concern, while a program should remain being interpreted when the programmer requires all the flexibility provided by the CX features.\n\n      </p>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Compiled and interpreted end -->\n\n\n<!-- Strict typing system start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"strict-typing-system\">Strict typing system</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n      <p>\n      There is no implicit casting in CX. Because of this, multiple versions for each of the primitive types are defined in the core module. For example, four native functions for addition exist: addI32, addI64, addF32, and addF64.\n      </p>\n      <p>\n      The parser attaches a default type to data it finds in the source code: if an integer is read, its default type is i32 or 32 bit integer; and if a float is read, its default type is f32 or 32 bit float. There is no ambiguity with other data read by the parser: true and false are always booleans; a series of characters enclosed between double quotes are always strings; and array needs to indicate its type before the list of its elements, e.g., []i64{{'{'}} 1, 2, 3{{'}'}} .\n      </p>\n      <p>\n      For the cases where the programmer needs to explicitly cast a value of one type to another, the core module provides a number of cast functions to work with primitive types. For example, byteAToStr casts a byte array to a string, and i32ToF32 casts a 32 bit integer to a 32 bit float.\n      </p>\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Strict typing system end -->\n\n\n<!-- Affordances start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"affordances\">Affordances</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n      <p>\n        A programmer needs to make a plethora of decisions while constructing a program, e.g., how many parameters a function must receive, how many parameters it must return, what statements are needed to obtain the desired functionality, and what arguments need to be sent as parameters to the statement functions, among others. The affordance system in CX can be queried to obtain a list of the possible actions that can be applied to an element.\n      </p>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Affordances end -->\n\n\n<!-- Serialization start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"serialization\">Serialization</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n      <p>\n      A program in CX can be partially or fully serialized to a byte array. This serialization capability allows a program to create a program image (similar to system images), where the exact state at which the program was serialized is maintained. This means that a serialized program can be deserialized, and resume its execution later on. Serialiation can also be used to create backups.\n      </p>\n      <p>\n      A CX program can leverage its integrated features to create some interesting scenarios. For example, a program can be serialized to create a backup of itself, and start an evolutionary algorithm on one of its functions. If the evolutionary algorithm finds a function that performs better than the previous definition, one can keep this new version of the program. However, if the evolutionary algorithm performed badly, the program can be restored to the saved backup. All of these tasks can be automated.\n      </p>\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Serialization end -->\n\n\n\n\n\n<!-- Stepping start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"stepping\">Stepping</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n      <p></p>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Stepping end -->\n\n\n\n\n\n<!-- Interactive evaluation start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"interactive-evaluation\">Interactive evaluation</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n      <p>\n      The affordance system and meta-programming functions in CX allow the flexibility of changing the program’s structure in a supervised manner. However, affordances can still be automated by having a function that selects the index of the affordance to be applied.\n      </p>\n      <p>\n        <i>evolve</i> is a native function that constructs user-defined functions by using random affordances.\n      </p>\n      <p>\n        <i>evolve</i> follows the principles of evolutionary computation. In particular, evolve performs a technique called genetic programming. Genetic programming tries to find a combination of operators and arguments that will solve a problem. For example, you could instruct evolve to find a combination of operators that, when sent 10 as an argument, returns 20. This might sound trivial, but genetic programming and other evolutionary algorithms can solve very complicated problems.\n      </p>\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Interactive evaluation end -->\n\n\n\n\n\n\n<!-- Interactive debugging start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"interactive-debugging\">Interactive debugging</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n      <p>\n        A CX program will enter the REPL mode once an error has been found. This behaviour gives the programmer the opportunity to debug the program before attempting to resume its execution.\n      </p>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n"

/***/ }),

/***/ "../../../../../src/app/components/tutorial/tutorial.comonent.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return TutorialComponent; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__ = __webpack_require__("../../../platform-browser/@angular/platform-browser.es5.js");
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var TutorialComponent = (function () {
    function TutorialComponent(titleService) {
        this.titleService = titleService;
    }
    TutorialComponent.prototype.ngOnInit = function () {
        this.titleService.setTitle('Tutorial');
    };
    return TutorialComponent;
}());
TutorialComponent = __decorate([
    Object(__WEBPACK_IMPORTED_MODULE_0__angular_core__["o" /* Component */])({
        selector: 'app-tutorial',
        template: __webpack_require__("../../../../../src/app/components/tutorial/tutorial.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/tutorial/tutorial.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["b" /* Title */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser__["b" /* Title */]) === "function" && _a || Object])
], TutorialComponent);

var _a;
//# sourceMappingURL=tutorial.comonent.js.map

/***/ }),

/***/ "../../../../../src/app/components/tutorial/tutorial.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/tutorial/tutorial.component.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"page-header\" id=\"banner\">\n  <div class=\"row\">\n    <div class=\"col-lg-6 col-md-6 col-sm-6\">\n      <h1>CX Tutorial</h1>\n      <p class=\"lead\">Contents</p>\n\n      <ul>\n        <li><a href=\"tutorial#introduction\">Introduction</a></li>\n        <ul>\n          <li><a href=\"tutorial#a-compiled-and-interpreted-language\" >A compiled and interpreted language</a></li>\n          <li><a href=\"tutorial#interactive-development\">Interactive development</a></li>\n          <li><a href=\"tutorial#meta-programming-commands\">Meta-programming commands</a></li>\n          <li><a href=\"tutorial#cxs-type-system\">CX's type system</a></li>\n          <li><a href=\"tutorial#cx-playground\">CX Playground</a></li>\n          <li><a href=\"tutorial#examples\">Examples</a></li>\n        </ul>\n        <li><a href=\"tutorial#program-architecture\">Program architecture</a></li>\n        <ul>\n          <li><a href=\"tutorial#comments\">Comments</a></li>\n          <li><a href=\"tutorial#packages\">Packages</a></li>\n          <li><a href=\"tutorial#global-definitions\">Global definitions</a></li>\n          <li><a href=\"tutorial#program-architecture-structs\">Structs</a></li>\n          <li><a href=\"tutorial#program-architecture-functions\">Functions</a></li>\n        </ul>\n        <li><a href=\"tutorial#flow-control-structures\">Flow-control structures</a></li>\n        <ul>\n          <li><a href=\"tutorial#if-and-if-else\">If and if/else</a></li>\n          <li><a href=\"tutorial#while\">While</a></li>\n          <li><a href=\"tutorial#go-to\">Go-to</a></li>\n        </ul>\n        <li><a href=\"tutorial#functions\">Functions</a></li>\n        <ul>\n          <li><a href=\"tutorial#multiple-return-values\">Multiple return values</a></li>\n          <li><a href=\"tutorial#named-result-parameters\">Named result parameters</a></li>\n          <li><a href=\"tutorial#expressions\">Expressions</a></li>\n        </ul>\n        <li><a href=\"tutorial#data\">Data</a></li>\n        <ul>\n          <li><a href=\"tutorial#definitions\">Definitions</a></li>\n          <li><a href=\"tutorial#data-structs\">Structs</a></li>\n          <li><a href=\"tutorial#arrays\">Arrays</a></li>\n        </ul>\n        <li><a href=\"tutorial#initialization\">Initialization</a></li>\n        <ul>\n          <li><a href=\"tutorial#variables\">Variables</a></li>\n        </ul>\n        <li><a href=\"tutorial#debugging\">Debugging</a></li>\n        <li><a href=\"tutorial#stepping\">Stepping</a></li>\n        <ul>\n          <li><a href=\"tutorial#program-halting\">Program halting</a></li>\n        </ul>\n        <li><a href=\"tutorial#affordances\">Affordances</a></li>\n        <li><a href=\"tutorial#standard-library\">Serialization</a></li>\n        <li><a href=\"tutorial#standard-library\">Standard library</a></li>\n        <ul>\n          <li><a href=\"tutorial#arithmetic-functions\">Arithmetic functions</a></li>\n          <li><a href=\"tutorial#flow-control-functions\">Flow-control functions</a></li>\n          <li><a href=\"tutorial#printing-functions\">Printing functions</a></li>\n          <li><a href=\"tutorial#identity-functions\">Identity functions</a></li>\n          <li><a href=\"tutorial#array-functions\">Array functions</a></li>\n          <li><a href=\"tutorial#casting-functions\">Casting functions</a></li>\n          <li><a href=\"tutorial#system-functions\">System functions</a></li>\n          <li><a href=\"tutorial#affordance-inference-functions\">Affordance inference functions</a></li>\n          <li><a href=\"tutorial#meta-programming-functions\">Meta-programming functions</a></li>\n          <li><a href=\"tutorial#random-numbers-functions\">Random numbers functions</a></li>\n          <li><a href=\"tutorial#logical-operators\">Logical operators</a></li>\n          <li><a href=\"tutorial#relational-operators\">Relational operators</a></li>\n          <li><a href=\"tutorial#bitwise-operators\">Bitwise operators</a></li>\n        </ul>\n      </ul>\n\n    </div>\n\n  </div>\n</div>\n\n<!-- Introduction start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"introduction\">Introduction</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n      <p>The <a href=\"/App/Specification\">CX specification</a> describes a set of data structures that define the different elements that a CX program can have, along with a set of functions/methods that manipulate these elements in order to execute a program. But these elements are part of a lower level of abstraction called CX Base. In principle, any language that follows the ideas behind CX Base can be considered a CX dialect.</p>\n\n      <p>This tutorial is based on an implementation of CX based on the lexicon, syntax and semantics of <a href=\"https://golang.org/\">Golang</a>. The repository of the whole CX project can be found <a href=\"https://github.com/skycoin/cx\">here</a>. The implementation described in this document can be found under the <i>src/cxgo/</i> directory.</p>\n\n      <h2 id=\"a-compiled-and-interpreted-language\">A compiled and interpreted language</h2>\n      <p>CX is designed to be both a compiled and an interpreted language. You can use the interpreted mode to interactively build and test a program, and when performance is an issue, the source code can be compiled. All of the features that can be used in the interpreted mode are accessible to its compiled counterpart, but it must be noted that CX is still a work in progress and this could change in the future.</p>\n\n      <h2 id=\"interactive-development\">Interactive development</h2>\n      <p>REPL. Add functions, expressions, other elements.</p>\n\n      <h2 id=\"meta-programming-commands\">Meta-programming commands</h2>\n      <p>Meta-programming commands. Selectors. Affordances. Stepping.</p>\n\n      <h2 id=\"cxs-type-system\">CX's type system</h2>\n\n      <p>\n        CX supports 13 primitive types:\n      </p>\n\n      <ul>\n        <li>str</li>\n        <li>bool</li>\n        <li>byte</li>\n        <li>i32</li>\n        <li>i64</li>\n        <li>f32</li>\n        <li>f64</li>\n        <li>[]bool</li>\n        <li>[]byte</li>\n        <li>[]i32</li>\n        <li>[]i64</li>\n        <li>[]f32</li>\n        <li>[]f64</li>\n      </ul>\n\n      <p>\n        In CX there is no type inference or implicit type casting. If a function requires a 32 bit integer, you can't send a 64 integer, or any of the two types of floats. For this reason, the <a href=\"#standard-library\">standard library</a> provides native functions that consider each of the primitive types when applicable. For example, there are four versions of the addition function: addI32, addI64, addF32, and addF64.\n      </p>\n\n      <p>\n        If, for example, you have an <i>i32</i> variable and you want to send it as an argument to a function which only accepts <i>f64</i>, you can use one of the casting functions to explicitly cast the variable to <i>f64</i>. To see the list of available casting functions, have a look at CX's <a href=\"#standard-library\">standard library</a> located at the bottom of this document. This is an example of type casting:\n      </p>\n\n      <pre><code class=\"go\">addF64(i32ToF64(3), f32ToF64(5.4))</code></pre>\n\n      <p>\n        Character strings in CX need to be enclosed in double quotes, and, unlike other languages such as JavaScript or Go, you can't use single quotes, backquotes or any other enclosing characters. Internally, strings are stored as byte arrays, and one can cast a string to a byte array, and viceversa:\n      </p>\n\n      <pre><code class=\"go\">byteToStr([]byte{{\"{\"}}0, 1, 2{{\"}\"}})\nstrToByte(\"hello\")</code></pre>\n\n      <p>\n        Booleans also have a special representation internally: they are actually 32 bit integers, limited to 1 (true) and 0 (false). Nevertheless, booleans can't be casted to numeric types. In CX, booleans are represented by the words <i>true</i> and <i>false</i>.\n      </p>\n\n      <p>\n        The programmer can create custom types by using CX structs, which are similar to structs in other programming languages. Structs in CX follow the syntax used in Go. Here is an example of a struct declaration:\n      </p>\n\n      <pre><code class=\"go\">type Point struct {{\"{\"}}\n\tname str\n\tx i32\n\ty i32\n}</code></pre>\n\n      <h2 id=\"cx-playground\">CX playground</h2>\n\n      <p>If you want to try CX programs, you can start doing it right away without downloading the actual interpreter/compiler. In the <a href=\"/\">home page</a> you can find a tool similar to <a href=\"https://play.golang.org/\">Go Playground</a>, where you can write and evaluate any CX program you want.</p>\n\n      <h2 id=\"examples\">Examples</h2>\n\n      <p>The full list of examples can be accessed by clicking <a href=\"/App/Examples\">here</a>, but one example is shown below, so you don't have to leave this webpage and for you to check how CX's syntax is like.</p>\n\n      <pre><code class=\"go\">package main\n\nfunc main () (out str) {{\"{\"}}\n     printStr(\"Hello World!\")\n{{\"}\"}}</code></pre>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Introduction end -->\n\n\n\n\n<!-- Program architecture start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"program-architecture\">Program architecture</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n      <p>The following sections describe the different elements that can be present in a CX program, and how they interact with each other.</p>\n\n      <h2 id=\"comments\">Comments</h2>\n      <p>There are two kinds of comments in CX: one for commenting out a single line, and another one to comment out a block of lines. As in other programming languages, commented lines are not evaluated by the interpreted or the compiler.</p>\n      <p>Any sequence of characters after two slashes (//) and until a newline character is found will be commented out, i.e., // comments out a single line.</p>\n      <p>Any text enclosed between a slash and an asterisk (/*) and an asterisk and a slash (*/) will be commented out, i.e., /* ... */ can comment out any number of lines.</p>\n\n      <pre><code class=\"go\">package main\n\nfunc main () (out i32) {{\"{\"}}\n\t// This won't be evaluated\n\t//out := divI32(5, 0) // This won't either\n\n\tout := divI32(10, 5) // The program will return 2\n\n\t/*\n        Comment block\n\tout := subF32(3.33, 1.11)\n\tout := subF32(3.33, 1.11)\n\tout := subF32(3.33, 1.11)\n\tout := subF32(3.33, 1.11)\n        */\n{{\"}\"}}</code></pre>\n\n      <h2 id=\"packages\">Packages</h2>\n      <p>Packages are a way to create groups of functions and definitions, which don't enter in conflict with other functions and definitions from other packages.</p>\n\n      <pre><code class=\"go\">package Math\n\nvar PI f32 = 3.14159\n\nfunc Square (num f32) (out f32) {{\"{\"}}\n\tmulF32(num, num)\n{{\"}\"}}\n\npackage Stat\n\nfunc Mean (vals []f32) (mean f32) {{\"{\"}}\n\tif eqI32(lenF32A(vals), 0) {{\"{\"}}\n\t\tprintStr(\"error?\")\n\t\thalt(\"Stat.Mean: division by 0\")\n\t{{\"}\"}}\n\tvar sum f32 = 0.0\n\tvar counter i32 = 0\n\twhile ltI32(counter, lenF32A(vals)) {{\"{\"}}\n\t\tsum = addF32(sum, readF32A(vals, counter))\n\t\tcounter = addI32(counter, 1)\n\t{{\"}\"}}\n\tmean = divF32(sum, i32ToF32(lenF32A(vals)))\n{{\"}\"}}\n\nfunc Variance (vals []f32) (variance f32) {{\"{\"}}\n\tmean := Mean(vals)\n\tvar sum f32 = 0.0\n\tvar counter i32 = 0\n\twhile ltI32(counter, lenF32A(vals)) {{\"{\"}}\n\t\tsum = Math.Square(subF32(readF32A(vals, counter), mean))\n\t\tcounter = addI32(counter, 1)\n\t{{\"}\"}}\n\tvariance = divF32(sum, i32ToF32(lenF32A(vals)))\n{{\"}\"}}\n        </code></pre>\n\n      <p>In the example above, two packages are created: one for math functions and definitions (PI and Square), and another for statistical functions (Mean and Variance). Variance calls a function from the Math package: Square. In order to specify that we want to call the function Square from the package Math, we need to use the full name of that function, which is Math.Square in this case. If we wanted to reference to the PI definition in the Math package, we'd use Math.PI.</p>\n\n      <h2 id=\"global-definitions\">Global definitions</h2>\n\n      <p>Global definitions are variables that are declared, and possibly initialized, outside of any function or struct. For example:</p>\n\n      <pre><code class=\"go\">package main\nvar key str = \"P6!7^P1Mme)LP+zcE=pH ^z_eg[3$OZY^rRg+D7R:-~7C/Db{{\"}\"}}w@Q&WX3vGZJXvVJ\"\nvar globalCounter i32\n\nfunc main () (out i32) {{\"{\"}}\n\tvar localCounter i32\n\tout := addI32(localCounter, globalCounter)\n{{\"}\"}}</code></pre>\n\n      <p>Any global or local variable which is not explicitly initialized will be automatically initialized to its corresponding zero-value. For example, 32 bit integer variables are initialized to 0, 32 bit floats are initialized to 0.0, strings to the empty character string (\"\"), booleans to false, etc. Other packages can access to global definitions too by using the definition's full name (e.g., Math.Sqrt. Read more about <a href=\"#packages\">packages</a> to know about accessing definitions and functions from other packages).</p>\n\n      <h2 id=\"program-architecture-structs\">Structs</h2>\n\n      <p>Structs, as in other programming laguages such as C or Golang, are used to create user-defined types. For example:</p>\n\n      <pre><code class=\"go\">type Student struct {{\"{\"}}\n\tname str\n\tcontrolNumber i32\n\tage i32\n\taddress str\n\tsmoker bool\n{{\"}\"}}</code></pre>\n\n      <p>The example above illustrates how a struct can be defined in CX. The keyword \"type\" starts a struct definition, and an identifier for that struct and the keyword \"struct\" need to follow. The group of variables or fields that will define the user-defined type are enclosed by braces. In this case, a Student type is defined, and has the fields: name, controlNumber, age, address, and a field that indicates if the student is a smoker or not.</p>\n\n      <h2 id=\"program-architecture-functions\">Functions</h2>\n\n      <p>Functions are the most complex elements in CX, as they can contain many other different elements. First, let's check a function definition that does absolutely nothing:</p>\n\n      <pre><code class=\"go\">func nothing () () {{\"{\"}}}</code></pre>\n\n      <p>Every function definition needs to start with the keyword \"func.\" Then we need to give this function a name, which is \"nothing\" in this example. The two pairs of parentheses that follow are used to indicate how many parameters the function is going to receive and return, as well as their names and types. Finally, any expression or statement that we want the function to run needs to be enclosed in the braces. Let's now define a function which prints a name to the console:</p>\n\n      <pre><code class=\"go\">func printName (name str) () {{\"{\"}}\n\tprintStr(name)\n{{\"}\"}}</code></pre>\n\n      <p>In this case, the function receives a single parameter: name. This parameter needs to be of the string type (str). The one and only expression that is going to be executed when calling this function is printStr(name), which, well, prints to the console the name that was sent. If we want to indicate that a function will receive or return multiple parameters, these need to be separated by a comma. Let's check out another example:</p>\n\n      <pre><code class=\"go\">func multiReturn (num1 i32, num2 i32) (add i32, sub i32, mul i32, div i32) {{\"{\"}}\n\tadd := addI32(num1, num2)\n\tsub := subI32(num1, num2)\n\tmul := mulI32(num1, num2)\n\tdiv := divI32(num1, num2)\n{{\"}\"}}</code></pre>\n\n      <p>The multiReturn function (extracted from the <a href=\"/App/Examples\">examples</a> section) illustrates how a multiple input and multiple output function can be defined. This function returns the sum, subtraction, multiplication and division of two numbers.</p>\n\n      <p>To learn more about the elements of a function, have a look at the CX <a href=\"/App/Examples\">examples</a> and the control structures, initialization sections.</p>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Program architecture end -->\n\n\n\n\n\n\n<!-- Control structures start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"flow-control-structures\">Flow-control structures</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n      <p></p>\n\n      <h2 id=\"if-and-if-else\">If and if/else</h2>\n\n      <p>The if and the if/else statements work just like in other programming languages, and its syntax is similar to Go's syntax:</p>\n\n      <pre><code class=\"go\">if gtI32(5, 3) {{\"{\"}}\n\tprintStr(\"5 is greater than 3\")\n{{\"}\"}}</code></pre>\n\n      <p>The keyword \"if\" starts an if statement. A condition or predicate needs to follow, which will indicate CX if the following block of statements can be executed or not.</p>\n\n      <pre><code class=\"go\">if eqStr(\"password123\", \"password123\") {{\"{\"}}\n\tprintStr(\"Access granted\")\n{{\"}\"}} else {{\"{\"}}\n\tprintStr(\"Access denied\")\n{{\"}\"}}</code></pre>\n\n      <p>If an \"else\" keyword follows the first block of statements, you can add a second block of statements. The first block is executed if the predicate evaluates to true, and the second block is executed if the predicate evaluates to false.</p>\n\n      <h2 id=\"while\">While</h2>\n\n      <p>The while statement is used to create loops in CX. The \"while\" keyboard starts a while statement, and a condition or predicate needs to follow it. If the predicate evaluates to true, the block of statements enclosed by braces is executed. After executing this block of statements, the predicate is re-evaluated. If the predicate evaluates again to true, the block of statements will be executed again. This process will be repeated until the predicate evaluates to false.</p>\n\n      <p>The example below shows a simple while loop that prints the numbers from 0 to 9 to the console.</p>\n\n      <pre><code class=\"go\">while ltI32(counter, 10) {{\"{\"}}\n\tprintI32(counter)\n\tcounter = addI32(counter, 1)\n{{\"}\"}}</code></pre>\n\n      <p>The while loop can also be used to iterate over arrays. An example of this process can be seen in the example below, where value in an array of 32 bit floats is accessed and summed to later calculate their mean.</p>\n\n      <pre><code class=\"go\">func Mean (vals []f32) (mean f32) {{\"{\"}}\n\tif eqI32(lenF32A(vals), 0) {{\"{\"}}\n\t\tprintStr(\"error?\")\n\t\thalt(\"Stat.Mean: division by 0\")\n\t{{\"}\"}}\n\tvar sum f32 = 0.0\n\tvar counter i32 = 0\n\twhile ltI32(counter, lenF32A(vals)) {{\"{\"}}\n\t\tsum = addF32(sum, readF32A(vals, counter))\n\t\tcounter = addI32(counter, 1)\n\t{{\"}\"}}\n\tmean = divF32(sum, i32ToF32(lenF32A(vals)))\n{{\"}\"}}</code></pre>\n\n      <h2 id=\"go-to\">Go-to</h2>\n\n      <p>All of the flow-control structures described in this section are internally parsed to structures that use the goTo native function. Actually, if we have a look at the <a href=\"standard-library\">standard library</a>, under the \"flow-control functions\", we'll see that there's only one function described: goTo.</p>\n\n      <p>The goTo function takes a predicate as its first parameter. If the predicate evaluates to true, the program will move forward or backward the number of lines defined by the second parameter. If the predicate evaluates to true, the number of lines to move forward or backward will be the number defined by the third parameter.</p>\n      <p>An example of its use is below. In this example, an implementation of an if statement is defined by the basicIf function.</p>\n\n      <pre><code class=\"go\">func basicIf (num i32) (num i32) {{\"{\"}}\n\tpred := gtI32(num, 0)\n\tgoTo(pred, 1, 3)\n\tprintStr(\"Greater than 0\")\n\tgoTo(true, 10, 0)\n\tprintStr(\"Less than 0\")\n{{\"}\"}}</code></pre>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Control structures end -->\n\n\n\n<!-- Initialization start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"initialization\">Initialization</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n      <p></p>\n\n      <h2 id=\"variables\">Variables</h2>\n\n      <p>Variables can be declared by using the \"var\" keyword inside a function. They are similar to definitions, with the difference that the scope of variables is limited to the function where they are declared. This is why they are sometimes referred to as locals, and definitions are sometimes referred to as globals.</p>\n\n      <pre><code class=\"go\">var foo str</code></pre>\n\n      <p>All variable declarations are implicitly initialized to the corresponding zero-value of the variable's type, if they are not explicitly initialized, of course. In the example above, the foo variable is not explicitly initialized, so its value after being declared is \"\" (an empty string). In order to do an explicit initialization, we can use the \"=\" keyword, followed by a value to be assigned to this variable:</p>\n\n      <pre><code class=\"go\">var foo str = \"Stay awhile and listen!\"</code></pre>\n\n      <p>Variables can also be declared and initialized with the values returned by expressions. In this case, the \"var\" keyword must not be used, and the \":=\" keyword has to be used instead of \"=\". Also, stating the type is not necessary and must not be added, as the variable will adopt the type of the value returned by the expression.</p>\n\n      <pre><code class=\"go\">foo := addI32(5, 10)</code></pre>\n\n      <p>In the example above, after evaluating the expression to the right of the \":=\" keyword, the foo variable will be of type i32, and will hold the value 15. The foo variable can change its value later on in the program by using the \"=\" keyword:</p>\n\n      <pre><code class=\"go\">foo = mulI32(10, 10)</code></pre>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Initialization end -->\n\n\n\n<!-- Debugging start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"debugging\">Debugging</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n\n      <p>In its development stage, a program in CX should be built and debugged in using the REPL. In the REPL mode, the programmer can control the execution of the program, interactively add and/or remove expressions, functions, structs, modules, etc. Check out the <a href=\"/App/Examples\">examples</a> section to learn more about meta-programming commands.</p>\n\n      <p>If a program is not executed in the REPL mode, and an error is encountered, the program will enter in the REPL mode to give the programmer or administrator the opportunity to fix the error. For example, the programmer can change the values of the local variables, function expressions can be altered, etc.</p>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Debugging end -->\n\n\n<!-- Standard library start -->\n<div class=\"bs-docs-section\">\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <div class=\"page-header\">\n        <h1 id=\"standard-library\">Standard library</h1>\n      </div>\n    </div>\n  </div>\n\n  <div class=\"row\">\n    <div class=\"col-lg-12\">\n      <!-- Content starts -->\n      <p>Every CX program has direct access to a package that defines the native functions. This package is called the \"core\" package, and it's not necessary to import it nor to append its name to the identifiers of the native functions, e.g., if you want to call addI32, it's not necessary to write core.addI32.</p>\n\n      <p>A description of each of the functions defined in the core package is provided below. These descriptions also include what parameters receive and return, as well as examples of their usage.</p>\n\n      <h2 id=\"arithmetic-functions\">Arithmetic functions</h2>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;addI32 (<i>number i32, number i32</i>) (<i>number i32</i>)</h4>\n          <h4>&#8594;addI64 (<i>number i64, number i64</i>) (<i>number i64</i>)</h4>\n          <h4>&#8594;addF32 (<i>number f32, number f32</i>) (<i>number f32</i>)</h4>\n          <h4>&#8594;addF64 (<i>number f64, number f64</i>) (<i>number f64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>This group of functions perform the addition between two numbers of the same type. A different version of the addition function is provided for each of the supported primitive types.</p>\n\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">apples := addI32(5, 3)\nviews := addI64(i32ToI64(1111), i32ToI64(3333))\nlitres := addF32(3.3, 4.4)\nerror := addF64(f32ToF64(0.00031331), f32ToF64(0.000025211))</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;subI32 (<i>number i32, number i32</i>) (<i>number i32</i>)</h4>\n          <h4>&#8594;subI64 (<i>number i64, number i64</i>) (<i>number i64</i>)</h4>\n          <h4>&#8594;subF32 (<i>number f32, number f32</i>) (<i>number f32</i>)</h4>\n          <h4>&#8594;subF64 (<i>number f64, number f64</i>) (<i>number f64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>This group of functions perform the subtraction between two numbers of the same type. A different version of the subtraction function is provided for each of the supported primitive types.</p>\n\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">apples := subI32(5, 3)\nviews := subI64(i32ToI64(3333), i32ToI64(1111))\nlitres := subF32(3.3, 4.4)\nerror := subF64(f32ToF64(0.00031331), f32ToF64(0.000025211))</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;mulI32 (<i>number i32, number i32</i>) (<i>number i32</i>)</h4>\n          <h4>&#8594;mulI64 (<i>number i64, number i64</i>) (<i>number i64</i>)</h4>\n          <h4>&#8594;mulF32 (<i>number f32, number f32</i>) (<i>number f32</i>)</h4>\n          <h4>&#8594;mulF64 (<i>number f64, number f64</i>) (<i>number f64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>This group of functions perform the multiplication between two numbers of the same type. A different version of the multiplication function is provided for each of the supported primitive types.</p>\n\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">apples := mulI32(5, 3)\nviews := mulI64(i32ToI64(3333), i32ToI64(1111))\nlitres := mulF32(3.3, 4.4)\nerror := mulF64(f32ToF64(0.00031331), f32ToF64(0.000025211))</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;divI32 (<i>number i32, number i32</i>) (<i>number i32</i>)</h4>\n          <h4>&#8594;divI64 (<i>number i64, number i64</i>) (<i>number i64</i>)</h4>\n          <h4>&#8594;divF32 (<i>number f32, number f32</i>) (<i>number f32</i>)</h4>\n          <h4>&#8594;divF64 (<i>number f64, number f64</i>) (<i>number f64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>This group of functions perform the division between two numbers of the same type. A different version of the division function is provided for each of the supported primitive types. If the denominator is equal to 0, a division by 0 error is raised.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">apples := divI32(5, 3)\nviews := divI64(i32ToI64(3333), i32ToI64(1111))\nlitres := divF32(3.3, 4.4)\nerror := divF64(f32ToF64(0.00031331), f32ToF64(0.000025211))</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;modI32 (<i>number i32, number i32</i>) (<i>number i32</i>)</h4>\n          <h4>&#8594;modI64 (<i>number i64, number i64</i>) (<i>number i64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>modI32 and modI64 perform the modulus operation between two numbers of the same type (i32 and i64, respectively). If the denominator is equal to 0, a division by 0 error is raised.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">remainder := modI32(5, 3)\nremainder := modI64(i32ToI64(3333), i32ToI64(1111))</code></pre>\n        </div>\n      </div>\n\n      <h2 id=\"flow-control-functions\">Flow-control functions</h2>\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;goTo (<i>predicate bool, thenLines i32, elseLines i32</i>) (<i>error bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>goTo evaluates a predicate, and if it evaluates to true, the program advances the number of lines defined by thenLines (second parameter), and if it evaluates to false, the program advances the number of lines defined by elseLines. thenLines and elseLines can be negative numbers.</p>\n\n          <p><b>Example:</b></p>\n          <pre><code class=\"go\">func basicIf (num i32) (num i32) {{\"{\"}}\n\tpred := gtI32(num, 0)\n\tgoTo(pred, 1, 3)\n\tprintStr(\"Greater than 0\")\n\tgoTo(true, 10, 0)\n\tprintStr(\"Less than 0\")\n{{\"}\"}}</code></pre>\n        </div>\n      </div>\n\n      <h2 id=\"printing-functions\">Printing functions</h2>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;printStr (<i>object str</i>) (<i>object str</i>)</h4>\n          <h4>&#8594;printBool (<i>object bool</i>) (<i>object bool</i>)</h4>\n          <h4>&#8594;printByte (<i>object byte</i>) (<i>object byte</i>)</h4>\n          <h4>&#8594;printI32 (<i>object i32</i>) (<i>object i32</i>)</h4>\n          <h4>&#8594;printI64 (<i>object i64</i>) (<i>object i64</i>)</h4>\n          <h4>&#8594;printF32 (<i>object f32</i>) (<i>object f32</i>)</h4>\n          <h4>&#8594;printF64 (<i>object f64</i>) (<i>object f64</i>)</h4>\n          <h4>&#8594;printBoolA (<i>object []bool</i>) (<i>object []bool</i>)</h4>\n          <h4>&#8594;printByteA (<i>object []byte</i>) (<i>object []byte</i>)</h4>\n          <h4>&#8594;printI32A (<i>object []i32</i>) (<i>object []i32</i>)</h4>\n          <h4>&#8594;printI64A (<i>object []i64</i>) (<i>object []i64</i>)</h4>\n          <h4>&#8594;printF32A (<i>object []f32</i>) (<i>object []f32</i>)</h4>\n          <h4>&#8594;printF64A (<i>object []f64</i>) (<i>object []f64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>This group of functions take a single argument as their input parameter, prints it the console, and then returns the argument as its output.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">printStr(\"greetings stranger!\")\nprintBool(false)\nprintByte(i32ToByte(50))\nprintI32(5)\nprintI64(i32ToI64(10))\nprintF32(3.14159)\nprintF64(f32ToF64(3.14159))\nprintBoolA([]bool{{\"{\"}}true, false, false{{\"}\"}})\nprintByteA([]byte{{\"{\"}}1, 2, 3, 4{{\"}\"}})\nprintI32A([]i32{{\"{\"}}1, 2, 3, 4{{\"}\"}})\nprintI64A([]i64{{\"{\"}}1, 2, 3, 4{{\"}\"}})\nprintF32A([]i64{{\"{\"}}1.5, 2.5, 3.5, 4.5{{\"}\"}})\nprintF64A([]i64{{\"{\"}}1.5, 2.5, 3.5, 4.5{{\"}\"}})</code></pre>\n        </div>\n      </div>\n\n      <h2 id=\"identity-functions\">Identity functions</h2>\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;idStr (<i>object str</i>) (<i>object str</i>)</h4>\n          <h4>&#8594;idBool (<i>object bool</i>) (<i>object bool</i>)</h4>\n          <h4>&#8594;idByte (<i>object byte</i>) (<i>object byte</i>)</h4>\n          <h4>&#8594;idI32 (<i>object i32</i>) (<i>object i32</i>)</h4>\n          <h4>&#8594;idI64 (<i>object i64</i>) (<i>object i64</i>)</h4>\n          <h4>&#8594;idF32 (<i>object f32</i>) (<i>object f32</i>)</h4>\n          <h4>&#8594;idF64 (<i>object f64</i>) (<i>object f64</i>)</h4>\n          <h4>&#8594;idBoolA (<i>object []bool</i>) (<i>object []bool</i>)</h4>\n          <h4>&#8594;idByteA (<i>object []byte</i>) (<i>object []byte</i>)</h4>\n          <h4>&#8594;idI32A (<i>object []i32</i>) (<i>object []i32</i>)</h4>\n          <h4>&#8594;idI64A (<i>object []i64</i>) (<i>object []i64</i>)</h4>\n          <h4>&#8594;idF32A (<i>object []f32</i>) (<i>object []f32</i>)</h4>\n          <h4>&#8594;idF64A (<i>object []f64</i>) (<i>object []f64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>This group of functions take a single argument as their input parameter, and then returns the argument as its output. These functions are the equivalent to the math function f(x) = x. These functions might not be of so much use in most programs, but CX uses them to parse programs and are provided to the user (like goTo).</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">printStr(\"greetings stranger!\")\nprintBool(false)\nprintByte(i32ToByte(50))\nprintI32(5)\nprintI64(i32ToI64(10))\nprintF32(3.14159)\nprintF64(f32ToF64(3.14159))\nprintBoolA([]bool{{\"{\"}}true, false, false{{\"}\"}})\nprintByteA([]byte{{\"{\"}}1, 2, 3, 4{{\"}\"}})\nprintI32A([]i32{{\"{\"}}1, 2, 3, 4{{\"}\"}})\nprintI64A([]i64{{\"{\"}}1, 2, 3, 4{{\"}\"}})\nprintF32A([]i64{{\"{\"}}1.5, 2.5, 3.5, 4.5{{\"}\"}})\nprintF64A([]i64{{\"{\"}}1.5, 2.5, 3.5, 4.5{{\"}\"}})</code></pre>\n        </div>\n      </div>\n\n      <h2 id=\"array-functions\">Array functions</h2>\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;readBoolA (<i>array []bool, index i32</i>) (<i>element bool</i>)</h4>\n          <h4>&#8594;readByteA (<i>array []byte, index i32</i>) (<i>element byte</i>)</h4>\n          <h4>&#8594;readI32A (<i>array []i32, index i32</i>) (<i>element i32</i>)</h4>\n          <h4>&#8594;readI64A (<i>array []i64, index i32</i>) (<i>element i64</i>)</h4>\n          <h4>&#8594;readF32A (<i>array []f32, index i32</i>) (<i>element f32</i>)</h4>\n          <h4>&#8594;readF64A (<i>array []f64, index i32</i>) (<i>element f64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>This group of functions read an element from the given array that corresponds to the element at the given index. If the index exceeds the length of the provided array, an error is raised. If the index is negative, an error is raised.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">readBoolA([]bool{{\"{\"}}true, false, true{{\"}\"}}, 0)\nreadByteA([]byte{{\"{\"}}1, 2, 3, 4{{\"}\"}}, 2)\nreadI32A([]i32{{\"{\"}}10, 20, 30, 40{{\"}\"}}, 3)\nreadI64A([]i64{{\"{\"}}10, 20, 30, 40{{\"}\"}}, 2)\nreadF32A([]f32{{\"{\"}}1.1, 2.2, 3.3, 4.4{{\"}\"}}, 1)\nreadF64A([]f64{{\"{\"}}1.1, 2.2, 3.3, 4.4{{\"}\"}}, 0)</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;writeBoolA (<i>array []bool, index i32, element bool</i>) (<i>array []bool</i>)</h4>\n          <h4>&#8594;writeByteA (<i>array []byte, index i32, element byte</i>) (<i>array []byte</i>)</h4>\n          <h4>&#8594;writeI32A (<i>array []i32, index i32, element i32</i>) (<i>array []i32</i>)</h4>\n          <h4>&#8594;writeI64A (<i>array []i64, index i32, element i64</i>) (<i>array []i64</i>)</h4>\n          <h4>&#8594;writeF32A (<i>array []f32, index i32, element f32</i>) (<i>array []f32</i>)</h4>\n          <h4>&#8594;writeF64A (<i>array []f64, index i32, element f64</i>) (<i>array []f64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>This group of functions write the given element to the given array, at the given index. If the index exceeds the length of the provided array, an error is raised. If the index is negative, an error is raised.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">writeBoolA([]bool{{\"{\"}}true, false, true{{\"}\"}}, 0)\nwriteByteA([]byte{{\"{\"}}1, 2, 3, 4{{\"}\"}}, 2, i32ToByte(10))\nwriteI32A([]i32{{\"{\"}}10, 20, 30, 40{{\"}\"}}, 3, 30)\nwriteI64A([]i64{{\"{\"}}10, 20, 30, 40{{\"}\"}}, 2, i32ToI64(50))\nwriteF32A([]f32{{\"{\"}}1.1, 2.2, 3.3, 4.4{{\"}\"}}, 1, 7.4)\nwriteF64A([]f64{{\"{\"}}1.1, 2.2, 3.3, 4.4{{\"}\"}}, 0, f32ToF64(10.10))</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;lenBoolA (<i>array []bool</i>) (<i>length i32</i>)</h4>\n          <h4>&#8594;lenByteA (<i>array []byte</i>) (<i>length i32</i>)</h4>\n          <h4>&#8594;lenI32A (<i>array []i32</i>) (<i>length i32</i>)</h4>\n          <h4>&#8594;lenI64A (<i>array []i64</i>) (<i>length i32</i>)</h4>\n          <h4>&#8594;lenF32A (<i>array []f32</i>) (<i>length i32</i>)</h4>\n          <h4>&#8594;lenF64A (<i>array []f64</i>) (<i>length i32</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>These functions calculate and return the length of the given array.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">lenBoolA([]bool{{\"{\"}}true, false, true{{\"}\"}})\nlenByteA([]byte{{\"{\"}}1, 2, 3, 4{{\"}\"}})\nlenI32A([]i32{{\"{\"}}10, 20, 30, 40{{\"}\"}})\nlenI64A([]i64{{\"{\"}}10, 20, 30, 40{{\"}\"}})\nlenF32A([]f32{{\"{\"}}1.1, 2.2, 3.3, 4.4{{\"}\"}})\nlenF64A([]f64{{\"{\"}}1.1, 2.2, 3.3, 4.4{{\"}\"}})</code></pre>\n        </div>\n      </div>\n\n      <h2 id=\"casting-functions\">Casting functions</h2>\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;byteAToStr (<i>array []byte</i>) (<i>string str</i>)</h4>\n          <h4>&#8594;strToByteA (<i>string str</i>) (<i>array []byte</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>A string is internally represented as an array of bytes. This implies that a byte array can be casted to a character string, and a character string can be casted to a byte array.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">byteAToStr([]byte{{\"{\"}}0, 1, 2{{\"}\"}})\nstrToByteA(\"hello\")</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;i64ToI32 (<i>number i64</i>) (<i>number i32</i>)</h4>\n          <h4>&#8594;f32ToI32 (<i>number f32</i>) (<i>number i32</i>)</h4>\n          <h4>&#8594;f64ToI32 (<i>number f64</i>) (<i>number i32</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>These functions cast their argument to an i32 number.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">i64ToI32(i32ToI64(5))\nf32ToI32(5.12)\nf64ToI32(f32ToF64(3.3))</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;i32ToI64 (<i>number i32</i>) (<i>number i64</i>)</h4>\n          <h4>&#8594;f32ToI64 (<i>number f32</i>) (<i>number i64</i>)</h4>\n          <h4>&#8594;f64ToI64 (<i>number f64</i>) (<i>number i64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>These functions cast their argument to an i64 number.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">i32ToI64(5)\nf32ToI64(5.12)\nf64ToI64(f32ToF64(3.3))</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;i32ToF32 (<i>number i32</i>) (<i>number f32</i>)</h4>\n          <h4>&#8594;i64ToF32 (<i>number i64</i>) (<i>number f32</i>)</h4>\n          <h4>&#8594;f64ToF32 (<i>number f64</i>) (<i>number f32</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>These functions cast their argument to an f32 number.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">i32ToF32(5)\ni64ToF32(i32ToI64(22))\nf64ToF32(f32ToF64(3.3))</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;i32ToF64 (<i>number i32</i>) (<i>number f64</i>)</h4>\n          <h4>&#8594;i64ToF64 (<i>number i64</i>) (<i>number f64</i>)</h4>\n          <h4>&#8594;f32ToF64 (<i>number f32</i>) (<i>number f64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>These functions cast their argument to an f64 number.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">i32ToF64(5)\ni64ToF64(i32ToI64(22))\nf32ToF64(3.3)</code></pre>\n        </div>\n      </div>\n\n\n\n\n\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;i64AToI32A (<i>array []i64</i>) (<i>array []i32</i>)</h4>\n          <h4>&#8594;f32AToI32A (<i>array []f32</i>) (<i>array []i32</i>)</h4>\n          <h4>&#8594;f64AToI32A (<i>array []f64</i>) (<i>array []i32</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>These functions cast their argument to an []i32 array.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">i64AToI32A([]i64{{\"{\"}}1, 2, 3{{\"}\"}})\nf32AToI32A([]f32{{\"{\"}}1.1, 2.2, 3.3{{\"}\"}})\nf64AToI32A([]f64{{\"{\"}}1.1, 2.2, 3.3{{\"}\"}})</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;i32AToI64A (<i>array []i32</i>) (<i>array []i64</i>)</h4>\n          <h4>&#8594;f32AToI64A (<i>array []f32</i>) (<i>array []i64</i>)</h4>\n          <h4>&#8594;f64AToI64A (<i>array []f64</i>) (<i>array []i64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>These functions cast their argument to an []i64 array.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">i32AToI64A([]i32{{\"{\"}}1, 2, 3{{\"}\"}})\nf32AToI64A([]f32{{\"{\"}}1.1, 2.2, 3.3{{\"}\"}})\nf64AToI64A([]f64{{\"{\"}}1.1, 2.2, 3.3{{\"}\"}})</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;i32AToF32A (<i>array []i32</i>) (<i>array f32</i>)</h4>\n          <h4>&#8594;i64AToF32A (<i>array []i64</i>) (<i>array f32</i>)</h4>\n          <h4>&#8594;f64AToF32A (<i>array []f64</i>) (<i>array f32</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>These functions cast their argument to an []f32 array.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">i32AToF32A([]i32{{\"{\"}}1, 2, 3{{\"}\"}})\ni64AToF32A([]i64{{\"{\"}}1, 2, 3{{\"}\"}})\nf64ToF32A([]f64{{\"{\"}}1.1, 2.2, 3.3{{\"}\"}})</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;i32AToF64A (<i>array []i32</i>) (<i>array []f64</i>)</h4>\n          <h4>&#8594;i64AToF64A (<i>array []i64</i>) (<i>array []f64</i>)</h4>\n          <h4>&#8594;f32AToF64A (<i>array []f32</i>) (<i>array []f64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>These functions cast their argument to an []f64 array.</p>\n\n          <p><b>Examples:</b></p>\n          <pre><code class=\"go\">i32AToF64A([]i32{{\"{\"}}1, 2, 3{{\"}\"}})\ni64AToF64A([]i64{{\"{\"}}1, 2, 3{{\"}\"}})\nf32ToF64A([]f32{{\"{\"}}1.1, 2.2, 3.3{{\"}\"}})</code></pre>\n        </div>\n      </div>\n\n      <h2 id=\"system-functions\">System functions</h2>\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;sleep (<i>milliseconds i32</i>) (<i>milliseconds i32</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Sleep halts a program for the amount of time defined by its argument. The amount of time is given in milliseconds.</p>\n\n          <p><b>Example:</b></p>\n          <pre><code class=\"go\">sleep(1000)</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;halt (<i>message str</i>) (<i>message str</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Halt is similar to sleep in that a program is paused temporarily. However, halt is usually used to inform the user that an error has been encountered. The program enters REPL mode after a call to halt has been encountered, so the programmer can start modifying the bugged program. Once the programmer has fixed the issues, a :step 0; meta-programming command can be issued to resume execution.</p>\n\n          <p><b>Example:</b></p>\n          <pre><code class=\"go\">halt(\"A score > 1000 has been reached.\")</code></pre>\n        </div>\n      </div>\n\n      <h2 id=\"affordance-inference-functions\">Affordance inference functions</h2>\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;setClauses (<i>clauses str</i>) (<i>clauses str</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Clauses are rules and/or facts used by the affordance system to determine if an affordance can be applied or not. setClauses receives a string containing prolog clauses. The programmer can provide the affordance system with any clauses, following any format, as long as the query (see setQuery below) can use them in combination with the defined objects (see addObject below).</p>\n\n          <p></p>\n\n          <p><b>Example:</b></p>\n          <pre><code class=\"go\">halt(\"A score > 1000 has been reached.\")</code></pre>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;setQuery (<i>query str</i>) (<i>query str</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;addObject (<i>object str</i>) (<i>object str</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;remObject (<i>object str</i>) (<i>object str</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;remObjects () (<i>error bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <h2 id=\"meta-programming-functions\">Random numbers functions</h2>\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;randI32 (<i>min i32, max i32</i>) (<i>number i32</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;randI64 (<i>min i64, max i64</i>) (<i>number i64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n\n      <h2 id=\"meta-programming-functions\">Meta-programming functions</h2>\n      <p>Selectors. Stepping. Debugging (dStack, dProgram, dState). Affordances. Removers. Prolog.</p>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;remExpr (<i>fnName str, num i32</i>) (<i>error bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;remArg (<i>fnName str</i>) (<i>error bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;addExpr (<i>fnName str</i>) (<i>error bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;exprAff (<i>filter str</i>) (<i>error bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;evolve (<i>fnName str, fnBag str, inputs []f64, outputs []f64, numberExprs i32, iterations i32, epsilon f64</i>) (<i>error f64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <!-- <h2 id=\"meta-programming-functions\">Meta-programming commands</h2>\n           <p>Selectors. Stepping. Debugging (dStack, dProgram, dState). Affordances. Removers. Prolog.</p> -->\n\n\n      <h2 id=\"logical-operators\">Logical operators</h2>\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;and (<i>premise bool, premise bool</i>) (<i>conclusion bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;or (<i>premise bool, premise bool</i>) (<i>conclusion bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;not (<i>premise bool</i>) (<i>conclusion bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <h2 id=\"relational-operators\">Relational operators</h2>\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;ltStr (<i>string str, string str</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;ltByte (<i>number byte, number byte</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;ltI32 (<i>number i32, number i32</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;ltI64 (<i>number i64, number i64</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;ltF32 (<i>number f32, number f32</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;ltF64 (<i>number f64, number f64</i>) (<i>conclusion bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;gtStr (<i>string str, string str</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;gtByte (<i>number byte, number byte</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;gtI32 (<i>number i32, number i32</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;gtI64 (<i>number i64, number i64</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;gtF32 (<i>number f32, number f32</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;gtF64 (<i>number f64, number f64</i>) (<i>conclusion bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;eqStr (<i>string str, string str</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;eqByte (<i>number byte, number byte</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;eqI32 (<i>number i32, number i32</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;eqI64 (<i>number i64, number i64</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;eqF32 (<i>number f32, number f32</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;eqF64 (<i>number f64, number f64</i>) (<i>conclusion bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;lteqStr (<i>string str, string str</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;lteqByte (<i>number byte, number byte</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;lteqI32 (<i>number i32, number i32</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;lteqI64 (<i>number i64, number i64</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;lteqF32 (<i>number f32, number f32</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;lteqF64 (<i>number f64, number f64</i>) (<i>conclusion bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;gteqStr (<i>string str, string str</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;gteqByte (<i>number byte, number byte</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;gteqI32 (<i>number i32, number i32</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;gteqI64 (<i>number i64, number i64</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;gteqF32 (<i>number f32, number f32</i>) (<i>conclusion bool</i>)</h4>\n          <h4>&#8594;gteqF64 (<i>number f64, number f64</i>) (<i>conclusion bool</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <h2 id=\"bitwise-operators\">Bitwise operators</h2>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;andI32 (<i>number i32, number i32</i>) (<i>number i32</i>)</h4>\n          <h4>&#8594;andI64 (<i>number i64, number i64</i>) (<i>number i64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;orI32 (<i>number i32, number i32</i>) (<i>number i32</i>)</h4>\n          <h4>&#8594;orI64 (<i>number i64, number i64</i>) (<i>number i64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;xorI32 (<i>number i32, number i32</i>) (<i>number i32</i>)</h4>\n          <h4>&#8594;xorI64 (<i>number i64, number i64</i>) (<i>number i64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <div class=\"panel panel-default\">\n        <div class=\"panel-heading\">\n          <h4>&#8594;andNotI32 (<i>number i32, number i32</i>) (<i>number i32</i>)</h4>\n          <h4>&#8594;andNotI64 (<i>number i64, number i64</i>) (<i>number i64</i>)</h4>\n        </div>\n        <div class=\"panel-body\">\n          <p>Lorem ipsum dolor sit amet</p>\n        </div>\n      </div>\n\n      <!-- Content ends -->\n    </div>\n  </div>\n\n</div>\n<!-- Standard library end -->\n\n"

/***/ }),

/***/ "../../../../../src/app/services/api.service.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return ApiService; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__("../../../http/@angular/http.es5.js");
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var ApiService = (function () {
    function ApiService(http) {
        this.http = http;
        this.api = '/';
    }
    ApiService.prototype.sendCode = function (code) {
        var url = this.api + 'eval';
        var data = {
            Code: code,
        };
        return this.http.post(url, data, new __WEBPACK_IMPORTED_MODULE_1__angular_http__["d" /* RequestOptions */]({
            headers: new __WEBPACK_IMPORTED_MODULE_1__angular_http__["a" /* Headers */]({ 'Content-Type': 'application/json' })
        }));
    };
    return ApiService;
}());
ApiService = __decorate([
    Object(__WEBPACK_IMPORTED_MODULE_0__angular_core__["C" /* Injectable */])(),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */]) === "function" && _a || Object])
], ApiService);

var _a;
//# sourceMappingURL=api.service.js.map

/***/ }),

/***/ "../../../../../src/environments/environment.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return environment; });
// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.
// The file contents for the current environment will overwrite these during build.
var environment = {
    production: false
};
//# sourceMappingURL=environment.js.map

/***/ }),

/***/ "../../../../../src/main.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__ = __webpack_require__("../../../platform-browser-dynamic/@angular/platform-browser-dynamic.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__app_app_module__ = __webpack_require__("../../../../../src/app/app.module.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__environments_environment__ = __webpack_require__("../../../../../src/environments/environment.ts");




if (__WEBPACK_IMPORTED_MODULE_3__environments_environment__["a" /* environment */].production) {
    Object(__WEBPACK_IMPORTED_MODULE_0__angular_core__["_23" /* enableProdMode */])();
}
Object(__WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__["a" /* platformBrowserDynamic */])().bootstrapModule(__WEBPACK_IMPORTED_MODULE_2__app_app_module__["a" /* AppModule */])
    .catch(function (err) { return console.log(err); });
//# sourceMappingURL=main.js.map

/***/ }),

/***/ 0:
/***/ (function(module, exports, __webpack_require__) {

module.exports = __webpack_require__("../../../../../src/main.ts");


/***/ })

},[0]);
//# sourceMappingURL=main.bundle.js.map

(window.webpackJsonp=window.webpackJsonp||[]).push([["chunk-65fec4c2"],{"02db":function(t,a,e){"use strict";(function(t){e("99af"),e("b0c0");var s=e("5f3f");a.a={data:function(){return{rememberPaw:"",CaptchaUrl:"",UNplace:"请输入用户名",PWplace:"请输入密码",verifyPWplace:"请输入确认密码",verPlace:"请输入验证码",loginForm:{name:"",password:"",captcha_id:"",verifyPassword:"",captcha_code:""},rules:{name:[{required:!0,message:"账号不能为空"}],password:[{required:!0,message:"密码不能为空"}],verifyPassword:[{required:!0,message:"密码不能为空"}],captcha_code:[{required:!0,message:"验证码不能为空"}]},carousel:[{img:e("2008"),path:"#"},{img:e("2008"),path:"#"},{img:e("2008"),path:"#"}]}},mounted:function(){this.getCaptchaId()},methods:{getCaptchaId:function(){var t=this;Object(s.a)().then((function(a){t.loginForm.captcha_id=a.id,t.CaptchaUrl="".concat(location.origin,"/admin/captcha/").concat(a.id,".png"),t.CaptchaUrl}))},thisYear:function(){return(new Date).getFullYear()},submit:function(a){var e=this,r=this;this.$refs[a].validate((function(a,o){if(o.name){var i=o.name[0].message;r.UNplace=i,t(".UNplace")[0].classList.add("err")}if(o.password){var c=o.password[0].message;r.PWplace=c,t(".PWplace")[0].classList.add("err")}if(o.captcha_code){var n=o.captcha_code[0].message;r.verPlace=n,t(".verPlace")[0].classList.add("err")}if(o.verifyPassword){var l=o.verifyPassword[0].message;r.verifyPWplace=l,t(".VPplace")[0].classList.add("err")}else if(r.loginForm.verifyPassword!=r.loginForm.password)return r.loginForm.verifyPassword="",r.verifyPWplace="两次密码输入不一致，请重新输入",void t(".VPplace")[0].classList.add("err");a&&Object(s.d)(e.loginForm).then((function(t){e.$message({message:"注册成功",type:"success"}),e.$router.push("/login")})).catch((function(t){}))}))}}}}).call(this,e("1157"))},2008:function(t,a,e){t.exports=e.p+"static/img/login1.c941d485.png"},5012:function(t,a,e){"use strict";e.r(a);var s=e("02db").a,r=e("2877"),o=Object(r.a)(s,(function(){var t=this,a=t.$createElement,e=t._self._c||a;return e("div",{staticClass:"container_login"},[e("div",{staticClass:"container_body"},[e("el-row",[e("el-col",{staticClass:"slideshow",attrs:{xs:12,sm:12,md:12,lg:12}},[e("div",{staticClass:"carousel"},[e("el-carousel",{attrs:{trigger:"click",height:"580px"}},t._l(t.carousel,(function(t,a){return e("el-carousel-item",{key:a},[e("a",{attrs:{href:t.path,target:"_blank"}},[e("img",{attrs:{src:t.img}})])])})),1)],1)]),e("el-col",{attrs:{xs:24,sm:12,md:12,lg:12}},[e("div",{staticClass:"loginForm"},[e("div",{staticClass:"form"},[e("div",{staticClass:"title"},[t._v("EasyRTC-SFU注册")]),e("div",{staticClass:"body"},[e("el-form",{ref:"loginForm",attrs:{model:t.loginForm,rules:t.rules,"show-message":!1}},[e("el-form-item",{attrs:{prop:"name"}},[e("div",{staticClass:"username"},[e("i",{staticClass:"iconfont iconadmin icon"}),e("input",{directives:[{name:"model",rawName:"v-model",value:t.loginForm.name,expression:"loginForm.name"}],staticClass:"formInput UNplace",attrs:{type:"text",placeholder:t.UNplace},domProps:{value:t.loginForm.name},on:{input:function(a){a.target.composing||t.$set(t.loginForm,"name",a.target.value)}}})])]),e("el-form-item",{attrs:{prop:"password"}},[e("div",{staticClass:"password"},[e("i",{staticClass:"iconfont iconpassword icon"}),e("input",{directives:[{name:"model",rawName:"v-model",value:t.loginForm.password,expression:"loginForm.password"}],staticClass:"formInput PWplace",attrs:{type:"password",placeholder:t.PWplace},domProps:{value:t.loginForm.password},on:{input:function(a){a.target.composing||t.$set(t.loginForm,"password",a.target.value)}}})])]),e("el-form-item",{attrs:{prop:"verifyPassword"}},[e("div",{staticClass:"password"},[e("i",{staticClass:"iconfont iconpassword icon"}),e("input",{directives:[{name:"model",rawName:"v-model",value:t.loginForm.verifyPassword,expression:"loginForm.verifyPassword"}],staticClass:"formInput VPplace",attrs:{type:"password",placeholder:t.verifyPWplace},domProps:{value:t.loginForm.verifyPassword},on:{input:function(a){a.target.composing||t.$set(t.loginForm,"verifyPassword",a.target.value)}}})])]),e("el-form-item",{attrs:{prop:"captcha_code"}},[e("div",{staticClass:"verification"},[e("input",{directives:[{name:"model",rawName:"v-model",value:t.loginForm.captcha_code,expression:"loginForm.captcha_code"}],staticClass:"formInput verPlace",attrs:{type:"text",placeholder:t.verPlace},domProps:{value:t.loginForm.captcha_code},on:{input:function(a){a.target.composing||t.$set(t.loginForm,"captcha_code",a.target.value)}}}),e("span",{staticClass:"verification_img",attrs:{title:"点击刷新"},on:{click:t.getCaptchaId}},[e("img",{attrs:{src:t.CaptchaUrl}})])])]),e("el-form-item",[e("div",{staticClass:"submit",on:{click:function(a){return t.submit("loginForm")}}},[t._v("注册")]),e("div",{staticClass:"loginTo"},[e("span",[t._v("已有账号?")]),e("router-link",{staticClass:"loginLink",attrs:{to:"/login"}},[t._v("请登录")])],1)])],1)],1)])])])],1)],1),e("div",{staticClass:"container_footer"},[e("span",{staticStyle:{color:"#808080"}},[t._v(" Copyright © 2014-"+t._s(t.thisYear())+" "),t._m(0),t._v(" All rights reserved ")])])])}),[function(){var t=this,a=t.$createElement,e=t._self._c||a;return e("a",{staticStyle:{color:"#2a88d7"},attrs:{href:"http://www.tsingsee.com/",target:"_target"}},[e("span",{staticStyle:{width:"70px",height:"16px",position:"relative",overflow:"hidden",display:"inline-block"}},[e("i",{staticClass:"iconfont iconqingxiLOGO",staticStyle:{"font-size":"70px",position:"absolute",top:"-15px",left:"0",color:"#2a88d7"}})]),t._v(" .com ")])}],!1,null,null,null);a.default=o.exports},"99af":function(t,a,e){"use strict";var s=e("23e7"),r=e("d039"),o=e("e8b5"),i=e("861d"),c=e("7b0b"),n=e("50c4"),l=e("8418"),d=e("65f0"),p=e("1dde"),m=e("b622"),u=e("2d00"),f=m("isConcatSpreadable"),g=9007199254740991,v="Maximum allowed index exceeded",h=u>=51||!r((function(){var t=[];return t[f]=!1,t.concat()[0]!==t})),w=p("concat"),y=function(t){if(!i(t))return!1;var a=t[f];return void 0!==a?!!a:o(t)};s({target:"Array",proto:!0,forced:!h||!w},{concat:function(t){var a,e,s,r,o,i=c(this),p=d(i,0),m=0;for(a=-1,s=arguments.length;a<s;a++)if(y(o=-1===a?i:arguments[a])){if(m+(r=n(o.length))>g)throw TypeError(v);for(e=0;e<r;e++,m++)e in o&&l(p,m,o[e])}else{if(m>=g)throw TypeError(v);l(p,m++,o)}return p.length=m,p}})},b0c0:function(t,a,e){var s=e("83ab"),r=e("9bf2").f,o=Function.prototype,i=o.toString,c=/^\s*function ([^ (]*)/,n="name";s&&!(n in o)&&r(o,n,{configurable:!0,get:function(){try{return i.call(this).match(c)[1]}catch(t){return""}}})}}]);
//# sourceMappingURL=chunk-65fec4c2.90a75ade.js.map
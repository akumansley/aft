(window.webpackJsonp=window.webpackJsonp||[]).push([[6],{76:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return o})),n.d(t,"metadata",(function(){return c})),n.d(t,"toc",(function(){return u})),n.d(t,"default",(function(){return s}));var r=n(3),a=n(7),i=(n(0),n(94)),o={id:"identity",title:"Identity"},c={unversionedId:"overview/identity",id:"overview/identity",isDocsHomePage:!1,title:"Identity",description:"Aft has a customizeable login system.",source:"@site/docs/overview/identity.md",slug:"/overview/identity",permalink:"/aft/overview/identity",version:"current",sidebar:"main",previous:{title:"Access",permalink:"/aft/overview/access"},next:{title:"Records",permalink:"/aft/overview/records"}},u=[],l={toc:u};function s(e){var t=e.components,n=Object(a.a)(e,["components"]);return Object(i.b)("wrapper",Object(r.a)({},l,n,{components:t,mdxType:"MDXLayout"}),Object(i.b)("p",null,"Aft has a customizeable login system."),Object(i.b)("p",null,"Out of the box, Aft has two ",Object(i.b)("a",Object(r.a)({parentName:"p"},{href:"rpcs"}),"RPCs"),": ",Object(i.b)("inlineCode",{parentName:"p"},"login")," and ",Object(i.b)("inlineCode",{parentName:"p"},"signup"),"."),Object(i.b)("p",null,"The code for each is very short. Here's login:"),Object(i.b)("pre",null,Object(i.b)("code",Object(r.a)({parentName:"pre"},{className:"language-python"}),'loginUnsuccessful = {"code": "login-error", "message": "login unsuccessful"}\n\ndef main(aft, args):\n    user = aft.api.findOne("user", {"where": {"email": args["email"]}})\n    if not user:\n        return loginUnsuccessful\n    if user.password == args["password"]:\n        aft.auth.authenticateAs(user.id)\n        return user\n    else:\n        return loginUnsuccessful\n')),Object(i.b)("p",null,"There are two methods on the ",Object(i.b)("inlineCode",{parentName:"p"},"aft.auth")," object: "),Object(i.b)("ol",null,Object(i.b)("li",{parentName:"ol"},Object(i.b)("inlineCode",{parentName:"li"},"authenticateAs(user_id)"),", which will generate an authentication token and inject it into the current request, as well as set it in a cookie"),Object(i.b)("li",{parentName:"ol"},Object(i.b)("inlineCode",{parentName:"li"},"user()"),", which will return the currently authenticated user object if there is one")),Object(i.b)("p",null,"The login and signup RPC are just regular code, that you may rewrite to provide any additional logic that's appropriate for your application."))}s.isMDXComponent=!0},94:function(e,t,n){"use strict";n.d(t,"a",(function(){return p})),n.d(t,"b",(function(){return d}));var r=n(0),a=n.n(r);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function c(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function u(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var l=a.a.createContext({}),s=function(e){var t=a.a.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):c(c({},t),e)),n},p=function(e){var t=s(e.components);return a.a.createElement(l.Provider,{value:t},e.children)},f={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},b=a.a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,i=e.originalType,o=e.parentName,l=u(e,["components","mdxType","originalType","parentName"]),p=s(n),b=r,d=p["".concat(o,".").concat(b)]||p[b]||f[b]||i;return n?a.a.createElement(d,c(c({ref:t},l),{},{components:n})):a.a.createElement(d,c({ref:t},l))}));function d(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var i=n.length,o=new Array(i);o[0]=b;var c={};for(var u in t)hasOwnProperty.call(t,u)&&(c[u]=t[u]);c.originalType=e,c.mdxType="string"==typeof e?e:r,o[1]=c;for(var l=2;l<i;l++)o[l]=n[l];return a.a.createElement.apply(null,o)}return a.a.createElement.apply(null,n)}b.displayName="MDXCreateElement"}}]);
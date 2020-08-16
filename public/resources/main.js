var headerExpanded = false;
var bar1 = document.getElementById("bar1");
var bar2 = document.getElementById("bar2");
var bar3 = document.getElementById("bar3");
var bar4 = document.getElementById("bar4");
var navA = document.querySelectorAll("header > nav > a");
var home = document.getElementById("home");
var homeText = document.querySelector("#home > span");
var search = document.querySelector("header > nav > form");
document.querySelector("header > nav > button").addEventListener("click", function() {
	headerExpanded = !headerExpanded;
	if(headerExpanded) {
		bar1.style.transform = "rotate(135deg)";
		bar2.style.transform = "rotate(135deg)";
		bar3.style.transform = "rotate(135deg)";
		bar4.style.transform = "rotate(-135deg)";
		bar4.style.background = "black";
		bar2.style.background = "transparent";
		bar3.style.background = "transparent";
		for(var i = 0; i < navA.length; i++) {
			navA[i].style.display = "block";
		}
		homeText.style.display = "block";
		home.style.position = "static";
		home.style.paddingTop = "10px";
		home.style.paddingBottom = "10px";
		home.style.paddingLeft = "10px";
		search.style.display = "block";
	} else {
		bar1.style.transform = "rotate(0deg)";
		bar2.style.transform = "rotate(0deg)";
		bar3.style.transform = "rotate(0deg)";
		bar4.style.transform = "rotate(0deg)";
		bar2.style.background = "black";
		bar3.style.background = "black";
		bar4.style.background = "transparent";
		search.style.display = "none";
		for(var i = 0; i < navA.length; i++) {
			navA[i].style.display = "none"
		}
		home.style.display = "block";
		homeText.style.display = "none";
	}
});

window.addEventListener("resize", function() {
	if(innerWidth > 720) {
		bar1.removeAttribute("style");
		bar2.removeAttribute("style");
		bar3.removeAttribute("style");
		bar4.removeAttribute("style");
		home.removeAttribute("style");
		homeText.removeAttribute("style");
		search.removeAttribute("style");
		for(var i = 0; i < navA.length; i++) {
			navA[i].removeAttribute("style");
		}
		headerExpanded = false;
	}
});

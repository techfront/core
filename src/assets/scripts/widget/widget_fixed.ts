/* TS for Fixed Widget */

namespace WidgetFixed {
    export class Main {
        fixedEl: HTMLElement;
        aboutEl: HTMLElement;
        scrollUpEl: HTMLElement;
        aboutIsClose: boolean;
        pageHeight: number;

        constructor() {
            this.fixedEl = <HTMLElement>Dom.First(".js-widget--fixed");
            this.aboutEl = <HTMLElement>Dom.First(".js-widget--fixed-about");
            this.scrollUpEl = <HTMLElement>Dom.First(".js-scroll-up");

            if (this.aboutEl !== undefined) {
                this.AboutCheckClose();
            }

            const body = document.body;
            const html = document.documentElement;

            window.onload = window.onscroll = () => {
                this.pageHeight = Math.max(body.scrollHeight, body.offsetHeight, html.clientHeight, html.scrollHeight, html.offsetHeight);
                this.Sticky();
                this.ScrollUpFadeIn();

                if (this.aboutEl !== undefined) {
                    this.AboutFadeOut();
                }
            };
        }

        Sticky(): void {
            const scroll = window.pageYOffset || document.documentElement.scrollTop;
            if ((this.pageHeight - scroll - document.documentElement.clientHeight) <= 150) {
                this.fixedEl.style.bottom = "38px";
                this.fixedEl.style.position = "absolute";
            } else {
                this.fixedEl.style.bottom = "";
                this.fixedEl.style.position = "";
            }
        }

        ScrollUpFadeIn(): void {
            if (window.pageYOffset <= 300) {
                this.scrollUpEl.style.opacity = "0.3";
                this.scrollUpEl.style.pointerEvents = "none";
            } else {
                this.scrollUpEl.style.opacity = ((window.pageYOffset) / 900).toString();
                this.scrollUpEl.style.pointerEvents = "";
            }
        }

        AboutFadeOut(): void {
            if (this.aboutIsClose !== true) {
                if (window.pageYOffset <= 520 ) {
                    this.aboutEl.style.opacity = ((500 -  window.pageYOffset) / 500).toString();
                    this.aboutEl.style.height = "120px";
                } else {
                    this.aboutEl.style.height = "0";
                }
            }
        }

        AboutCheckClose(): void {
            if (localStorage.getItem("techfront_.js-widget--fixed-about_close") === "true") {
                this.aboutEl.parentNode.removeChild(this.aboutEl);
                this.aboutIsClose = true;
            } else {
                this.aboutEl.style.display = "";
                this.aboutIsClose = false;
            }
        }
    }
}
// Package DOM provides functions to replace the use of jquery in 1.4KB of js - Ajax, Selectors, Event binding, ShowHide
// See http://youmightnotneedjquery.com/ for more if required
// Version 1.0.1

class Dom {

    static Ready(f: () => void): void {
        if (document.readyState !== "loading") {
            f();
        } else {
            document.addEventListener("DOMContentLoaded", f);
        }
    }

    static Exists(s: string): boolean {
        return (document.querySelector(s) !== null);
    }

    static All(s: string): NodeListOf < Element > {
        return document.querySelectorAll(s);
    }

    static Nearest(el: Element, s: string): NodeListOf < Element > | void {
        while (el !== undefined) {
            const nearest = el.querySelectorAll(s);
            if (nearest.length > 0) {
                return nearest;
            }

            el = < Element > el.parentNode;
        }
        return undefined;
    }

    static First(s: string): Element {
        return this.All(s)[0];
    }

    static Each(s: string, f: Function): void {
        const a = this.All(s);
        for (let i = 0; i < a.length; ++i) {
            f(a[i], i);
        }
    }

    static ForEach(a: Array < any > , f: Function): void {
        Array.prototype.forEach.call(a, f);
    }

    static Hide(s: string): void {
        this.Each(s, (el: HTMLElement, i: number) => {
            el.style.display = "none";
        });
    }

    static Show(s: string): void {
        this.Each(s, (el: HTMLElement, i: number) => {
            el.style.display = "";
        });
    }

    static ShowHide(s: string): void {
        this.Each(s, (el: HTMLElement, i: number) => {
            if (el.style.display !== "none") {
                this.Hide(s);
            } else {
                this.Show(s);
            }
        });
    }

    static On(s: string, type: string, callback: Function): void {
        this.Each(s, (element: Element, i: number) => {
            element.addEventListener(type, <EventListenerOrEventListenerObject>callback);
        });
    }

    static One(s: string, type: string, callback: Function): void {
        this.Each(s, (element: Element, i: number) => {
            element.addEventListener(type, <EventListenerOrEventListenerObject>function(e: Event) {
                e.target.removeEventListener(e.type, <EventListenerOrEventListenerObject>arguments.callee);
                return callback(e);
            });
        });
    }

    static Format(f: string): string {
        for (let i = 1; i < arguments.length; i++) {
            const regexp = new RegExp(`\\{${i - 1}\\}`, "gi");
            f = f.replace(regexp, arguments[i]);
        }
        return f;
    }

    static Post(u: string, d: any, fs: Function, fe: Function): void {
        const request = new XMLHttpRequest();
        request.onreadystatechange = () => {
            if (request.readyState === 4 && request.status === 200) {
                fs(request);
            } else {
                fe(request);
            }
        };
        request.open("POST", u, true);
        request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8");
        request.send(d);
    }

    static Get(u: string, fs: Function, fe: Function): void {
        const request: any = new XMLHttpRequest();
        request.open("GET", u, true);
        request.onload = () => {
            if (request.status >= 200 && request.status < 400) {
                fs(request);
            } else {
                fe();
            }
        };
        request.onerror = fe;
        request.send();
    }
}

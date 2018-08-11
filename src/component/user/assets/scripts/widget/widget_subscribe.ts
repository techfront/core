/* TS for Subscribe Widget */

namespace WidgetSubscribe {
    export class Main {
        element: HTMLElement;
        private notifier: Notifier;

        constructor() {
            this.element = <HTMLElement>Dom.First(".js-widget--subscribe");

            if (this.element !== undefined) {
                this.notifier = new Notifier["default"]({
                    theme: "techfront",
                    position: "top-right"
                });
                this.HandleSubmit();
                this.CheckClose();
            }
        }

        HandleSubmit() {
            Dom.On(".js-widget--subscribe-form", "submit", (e: Event) => {
                e.preventDefault();
                const form = <HTMLElement>e.target;
                this.FetchResponse(form);
            });
        }

        FetchResponse(form: HTMLElement) {
            const data = getFormData(form, { trim: true });

            let params = "";
            for (let i in data) {
                params += encodeURIComponent(i) + "=" + encodeURIComponent(data[i]) + "&";
            }
            window.fetch("/api/v3/post/user/subscribe", {
                method: "POST",
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8"
                },
                body: params
            })
            .then((response: any) => response.json())
            .then((response: any) => {
                if (typeof response._code !== "undefined") {
                    this.notifier.post(response._meta._notice, { delay: 6000, type: "error" });
                } else {
                    this.notifier.post(response._meta._notice, { delay: 6000, type: "success" });
                }
            });
        }

        CheckClose(): void {
            if (localStorage.getItem("techfront_.js-widget--subscribe_close") === "true") {
                this.element.parentNode.removeChild(this.element);
            } else {
                this.element.style.display = "";
            }
        }
    }
}
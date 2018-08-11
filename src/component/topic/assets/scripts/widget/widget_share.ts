/* TS for Share Widget */

namespace WidgetShare {
    export class Main {
        constructor() {
            this.ActivateShareLinks();
        }

        ActivateShareLinks(): void {
            Dom.On(".js-widget--share", "click", function(e: Event) {
                let width: number = 650,
                    height: number = 450;

                e.preventDefault();

                window.open(this.href, "Share Dialog", "menubar=no,toolbar=no,resizable=yes,scrollbars=yes,width=" + width + ",height=" + height + ",top=" + (screen.height / 2 - height / 2) + ",left=" + (screen.width / 2 - width / 2));
            });
        }
    }
}
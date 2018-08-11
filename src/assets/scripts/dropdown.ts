/* TS for Dropdowns */

namespace Dropdown {
    export class Main {
        constructor() {
            this.ActivateToggleLinks();
        }

        CloseByOutsideClick(element: Element): void {
            Dom.One("html", "click", (e: Event) => {
                element.classList.remove("dropdown__menu--show");
            });
        }

        ActivateToggleLinks() {
            Dom.On(".js-dropdown", "click", (e: Event) => {
                const target = <Element>e.currentTarget;
                const selector = target.getAttribute("data-menu");
                const menu = Dom.First(selector);
                menu.classList.toggle("dropdown__menu--show");
                e.preventDefault();
                if (menu.classList.contains("dropdown__menu--show") !== false) {
                    this.CloseByOutsideClick(menu);
                    e.stopPropagation();
                }
            });
        }
    }
}
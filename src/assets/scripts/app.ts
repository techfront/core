/* Basic TS */

namespace App {
    export class Main {
        constructor() {
            // Check Links
            this.CheckVoteLinks();

            // Activate Links
            this.ActivateVoteLinks();
            this.ActivateCloseLinks();
            this.ActivateShowLinks();
            this.ActivateMethodLinks();
            this.ActivateMessageLinks();
            this.ActivateModalLinks();

            // Toolbar
            this.ScrollUp();
            this.ScrollPrev();
            this.ScrollNext();
        }

        ActivateCloseLinks(): void {
            Dom.On(".js-close", "click", function() {
                const selector = this.getAttribute("data-close");
                const element = Dom.First(selector);
                element.parentNode.removeChild(element);
                localStorage.setItem("techfront_" + selector + "_close", "true");
            });
        }

        ActivateMethodLinks(): void {
            Dom.On("a[method='post'], a[method='delete']", "click", function(e: Event) {
                // Confirm action before delete
                if (this.getAttribute("method") === "delete") {
                    if (!confirm("Are you sure you want to delete this item, this action cannot be undone?")) {
                        return false;
                    }
                }

                // Collect the authenticity token from meta tags in header
                const meta = Dom.First("meta[name='authenticity_token']");
                if (meta === undefined) {
                    e.preventDefault();
                    return false;
                }
                const token = meta.getAttribute("content");

                // Perform a post to the specified url (href of link)
                const url = this.getAttribute("href");
                const redirect = this.getAttribute("data-redirect");
                const data = `authenticity_token=${token}`;

                Dom.Post(url, data, () => {
                    if (redirect !== null) {
                        // If we have a redirect, redirect to it after the link is clicked
                        (window as any).location = redirect;
                    } else {
                        // If no redirect supplied, we just reload the current screen
                        window.location.reload();
                    }
                }, () => {
                });

                e.preventDefault();
                return false;
            });


            Dom.On("a[method='back']", "click", (e: Event) => {
                history.back();
                e.preventDefault();
                return false;
            });

        }

        ActivateShowLinks(): void {
            Dom.On(".js-show", "click", function(e: Event) {
                const selector = this.getAttribute("data-show");

                Dom.Each(selector, (element: HTMLElement, i: number) => {
                    if (!element.className.match(/hidden/gi)) {
                        element.className = `${element.className}--hidden`;
                    } else {
                        element.className = element.className.replace(/--hidden/gi, "");
                    }
                });

                e.preventDefault();

                return false;
            });
        }

        CheckVoteLinks(): void {
            Dom.Each("[data-vote]", (element: HTMLElement, i: number) => {
                let vote = element.getAttribute("data-vote");
                let voteID = element.getAttribute("data-vote-id");

                if (localStorage.getItem("techfront_" + voteID + "_vote") === vote) {
                    let parentElement = <HTMLElement>element.parentNode;

                    if (!parentElement.className.match(/notvoted/gi)) {
                        parentElement.className = `${element.className}--notvoted`;
                    } else {
                        parentElement.className = parentElement.className.replace(/--notvoted/gi, "");
                    }
                }
            });
        }

        ActivateVoteLinks(): void {
            Dom.On(".js-vote", "click", function() {
                const vote = this.getAttribute("data-vote");
                const voteID = this.getAttribute("data-vote-id");
                localStorage.setItem("techfront_" + voteID + "_vote", vote);
            });
        }

        SmoothScroll(offset: number, duration: number, e: Event): void {
            e.preventDefault();

            let cosParameter = (window.scrollY - offset) / 2,
                scrollCount = 0,
                oldTimestamp = performance.now();

            const requestAnimationFrame = window.requestAnimationFrame ||
                (window as any).mozRequestAnimationFrame ||
                (window as any).webkitRequestAnimationFrame ||
                (window as any).msRequestAnimationFrame;

            function step(newTimestamp: number): void {
                scrollCount += Math.PI / (duration / (newTimestamp - oldTimestamp));
                if (scrollCount >= Math.PI) return;

                window.scrollTo(0, Math.round(offset + cosParameter + cosParameter * Math.cos(scrollCount)));
                oldTimestamp = newTimestamp;
                requestAnimationFrame(step);
            }

            requestAnimationFrame(step);
        }

        ScrollUp(): void {
            Dom.On(".js-scroll-up", "click", (e: Event) => {
                this.SmoothScroll(0, 1000, e);
            });
        }

        ScrollNext(): void {
            Dom.On(".js-control-next", "click", (e: Event) => {
                let offset = window.scrollY + 500;
                this.SmoothScroll(offset, 600, e);
            });
        }

        ScrollPrev(): void {
            Dom.On(".js-control-prev", "click", (e: Event) => {
                let offset = window.scrollY - 500;
                this.SmoothScroll(offset, 600, e);
            });
        }

        ActivateMessageLinks(): void {
            Dom.On(".js-message-close", "click", function() {
                const selector = this.getAttribute("data-message");
                const element = Dom.First(selector);
                element.parentNode.removeChild(element);
            });
        }

        ActivateModalLinks(): void {
            Dom.On(".js-modal-show", "click", function(e: Event) {

                // Скрываем другие окна
                Dom.Each(".modal", (element: HTMLElement, i: number) => {
                    element.classList.add("modal--hidden");
                });

                // Узнаем какое окно открыть
                const selector = this.getAttribute("data-modal");

                // Открываем окно
                Dom.Each(selector, (element: HTMLElement, i: number) => {
                    element.classList.remove("modal--hidden");
                });

                e.preventDefault();
                return false;
            });

            Dom.On(".js-modal-hide", "click", function(e: Event) {
                const selector = this.getAttribute("data-modal");

                Dom.Each(selector, (element: HTMLElement, i: number) => {
                    element.classList.add("modal--hidden");
                });

                e.preventDefault();
                return false;
            });
        }

    }
}
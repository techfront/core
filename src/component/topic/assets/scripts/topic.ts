/* TS for Topic Component */

namespace ComponentTopic {
    export interface FavoriteResponseMeta {
        _notice: string;
    }

    export interface FavoriteResponse {
        _code: string;
        _meta: FavoriteResponseMeta;
    }
    export class Main {
        private notifier: Notifier;

        constructor() {
           this.notifier = new Notifier["default"]({
                  theme: "techfront",
                  position: "top-right"
           });
            window.addEventListener("DOMContentLoaded", () => {
                this.SetBlockContainerHeight();
            });
            window.addEventListener("resize", () => {
                this.SetBlockContainerHeight();
            });
           this.HandleFavorite();
        }

        SetBlockContainerHeight(): void {
            const containerElement = <HTMLElement>Dom.First(".js-topic-index--block-container");
            if (containerElement !== undefined || typeof containerElement !== "undefined") {
                if (containerElement.querySelector(".block") !== null) {
                    containerElement.style.height = (containerElement.firstElementChild.clientHeight - 10).toString() + "px";

                    Dom.On(".js-topic-index--block-container .js-close", "click", () => {
                        containerElement.style.height = "";
                    });
                }
            }
        }

        async FetchCreateFavorite(userId, topicId: string): Promise<FavoriteResponse> {
           const response = await window.fetch("/api/v3/post/user/favorite", {
                method: "POST",
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8"
                },
                body: "id_user=" + userId + "&id_topic=" + topicId,
                credentials: "same-origin"
            });
            return await response.json();
        }

        async FetchDestroyFavorite(userId, topicId: string): Promise<FavoriteResponse> {
           const response = await window.fetch("/api/v3/post/user/unfavorite", {
                method: "POST",
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8"
                },
                body: "id_user=" + userId + "&id_topic=" + topicId,
                credentials: "same-origin"
            });
            return await response.json();
        }

        ToggleFavoriteLink(element: HTMLElement, action: string): void {
            if (action === "create") {
                element.innerHTML = "<i class=\"icon-bookmark\"></i> Удалить из избранного";
                element.className = "js-topic--favorite topic__meta-unfavorite";
                element.setAttribute("data-action", "destroy");
            } else if (action === "destroy") {
                element.innerHTML = "<i class=\"icon-bookmark\"></i> Прочитать позже";
                element.className = "js-topic--favorite topic__meta-favorite";
                element.setAttribute("data-action", "create");
            }
        }

        HandleFavorite() {
            Dom.On(".js-topic--favorite", "click", (e: Event) => {
                const target: HTMLElement = <HTMLElement>e.target;
                const action: string = target.getAttribute("data-action");
                const userId: string = target.getAttribute("data-id-user");
                const topicId: string = target.getAttribute("data-id-topic");
                let response: Promise<FavoriteResponse>;
                if (action === "create") {
                    response = this.FetchCreateFavorite(userId, topicId);
                } else if (action === "destroy") {
                    response = this.FetchDestroyFavorite(userId, topicId);
                }
                response.then(data => {
                    if (typeof data._code !== "undefined") {
                        this.notifier.post(data._meta._notice, { delay: 6000, type: "error" });
                    } else {
                        this.ToggleFavoriteLink(target, action);
                        this.notifier.post(data._meta._notice, { delay: 6000, type: "success" });
                    }
                }).catch(reason => {
                    console.log(reason.message);
                });
                e.preventDefault();
            });
        }
    }
}
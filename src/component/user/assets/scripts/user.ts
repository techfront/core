/* TS for User Component */

namespace ComponentUser {
    export class Main {
        activateMessage: HTMLElement;
        private notifier: Notifier;

        constructor() {
                this.activateMessage = <HTMLElement>Dom.First(".js-message--user-activate");
                if (this.activateMessage !== undefined) {
                     this.notifier = new Notifier["default"]({
                          theme: "techfront",
                          position: "top-right"
                     });
                     this.HandleActivate();
                }
                this.ContactFieldAdd();
                this.ContactFieldRemove();
                this.ContactFieldChangeType();
        }

        HandleActivate() {
            Dom.On(".js-message--user-activate-send", "click", (e: Event) => {
                const target: HTMLElement = <HTMLElement>e.target;
                const id = target.getAttribute("data-user-id");
                window.fetch("/api/v3/post/user/activate", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8"
                    },
                    body: "id_user=" + id,
                    credentials: "same-origin"
                })
                .then((response: any) => response.json())
                .then((response: any) => {
                    if (typeof response._code !== "undefined") {
                        this.notifier.post(response._meta._notice, { delay: 6000, type: "error" });
                    } else {
                        this.notifier.post(response._meta._notice, { delay: 6000, type: "success" });
                    }
                });
                e.preventDefault();
            });
        }

        ContactFieldChangeType() {
            Dom.On(".js-user--contact-field-change-type", "change", (e: Event) => {
                const select: HTMLSelectElement = <HTMLSelectElement>e.target;
                const selectValue = select.options[select.selectedIndex].value;
                const fieldId = select.getAttribute("data-field-id");
                const field = document.getElementById(fieldId);
                field.querySelector(".input").setAttribute("name", selectValue);
            });
        }

        ContactFieldRemove() {
            Dom.On(".js-user--contact-fields-remove", "click", function(e: Event) {
                const container = document.getElementById("user-contact-fields-list");
                const fieldId = this.getAttribute("data-field-id");
                const field = document.getElementById(fieldId);
                field.parentNode.removeChild(field);
                if (container.childElementCount < 8) {
                    const cfaTarget = document.querySelector(".js-user--contact-fields-add");
                    cfaTarget.classList.remove("user__contact-fields-add--lock");
                }
                e.preventDefault();
            });
        }

        ContactFieldAdd() {
            Dom.On(".js-user--contact-fields-add", "click", (e: Event) => {
                const target: HTMLElement = <HTMLElement>e.target;
                const container = document.getElementById("user-contact-fields-list");
                if (target.classList.contains("user__contact-fields-add--lock")) {
                    e.preventDefault();
                    return;
                }

                const fieldId = "user-contact-field-id-" + container.childElementCount;
                const field = document.createElement("div");
                field.id = fieldId;
                field.classList.add("user__contact-fields-item");
                field.innerHTML = `
                                 <select class="js-user--contact-field-change-type select" data-field-id="` + fieldId + `">
                                    <option value="user_contact_email" selected="selected">Почта</option>
                                    <option value="user_contact_phone">Телефон</option>      
                                    <option value="user_contact_website">Сайт</option>
                                     <option value="user_contact_github">Github</option>
                                    <option value="user_contact_telegram">Telegram</option>
                                    <option value="user_contact_whatsapp">WhatsApp</option>
                                    <option value="user_contact_viber">Viber</option>
                                    <option value="user_contact_skype">Skype</option>
                                    <option value="user_contact_vkontakte">ВКонтакте</option>
                                    <option value="user_contact_twitter">Twitter</option>
                                    <option value="user_contact_youtube">YouTube</option>
                                    <option value="user_contact_facebook">Facebook</option>
                                </select>
                                <input type="text" class="input" name="user_contact_email" value="">
                                <a href="#" class="js-user--contact-fields-remove user__contact-fields-remove" data-field-id="` + fieldId + `"><i class="icon-trash-empty"></i></a>`;
                container.appendChild(field);
                this.ContactFieldChangeType();
                this.ContactFieldRemove();
                if (container.childElementCount >= 8) {
                    target.classList.add("user__contact-fields-add--lock");
                }
                e.preventDefault();
            });
        }
    }
}
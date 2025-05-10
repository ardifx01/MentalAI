const optionsTanggal = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
};

function UpdatePercakapan(DataJson) {
    const chatDetailMessages = document.getElementById('chat-detail-messages');
    if(!DataJson) {
        chatDetailMessages.innerHTML = document.getElementById('empty-state-template').innerHTML;
    } else {
        chatDetailMessages.innerHTML = "";

        const JudulChat = document.getElementById("judul-chat");
        if(JudulChat) {
            JudulChat.innerText = DataJson.judul;

            if(location.pathname.includes("history")) {
                $(JudulChat).append(`<span class="urgency-badge urgency-${DataJson.urgency_level}" id="urgencylevel-chat">${DataJson.urgency_level}</span>`);
            }
        }

        const UrgencyLevelChat = document.getElementById("urgencylevel-chat");
        if(UrgencyLevelChat) {
            UrgencyLevelChat.textContent = DataJson.urgency_level;
            UrgencyLevelChat.classList.add("urgency-" + `${DataJson.urgency_level}`);
        }

        const WaktuChat = document.getElementById("waktu-chat");
        if(WaktuChat) {
            document.getElementById("waktu-chat").textContent = "Started: " + new Date(DataJson.created_at).toLocaleString('en-US', optionsTanggal);
        }

        if(location.pathname.includes("history")) {
            chatDetailMessages.innerHTML = `
            <div class="date-divider">
                <span>${new Date(DataJson.created_at).toLocaleString('en-US', optionsTanggal)}</span>
            </div>
            `;
        }

        if(location.pathname.includes("dashboard")) {
            chatDetailMessages.innerHTML = `
            <div class="message message-bot">
                <p class="mb-0">Hello! I'm your MindfulAI assistant. How are you feeling today?</p>
            </div>
            `;
        }

        for(const [_, v] of Object.entries(DataJson.omongan)) {
            const waktu = new Date(v.created_at).toLocaleTimeString('en-US', {
                hour: 'numeric',
                minute: '2-digit',
                hour12: true
            });                
            const messageClass = v.pengirim === "user" ? "message-user" : "message-bot";
            const messageTime = `<div class="message-time">${waktu}</div>`;
            const messageContent = `<p class="mb-0">${marked.parse(v.pesan)}</p>`;
            chatDetailMessages.innerHTML += `
                <div class="message ${messageClass}">
                    ${messageContent}
                    ${(location.pathname.includes("history") ? messageTime : "")}
                </div>
            `;
        }

        if(location.pathname.includes("history")) {
            chatDetailMessages.innerHTML += `
            <div class="date-divider">
                <span>End of Conversation</span>
            </div>
            `
        }
    }

    percakapan_id = DataJson.id;
    document.getElementById('chat').setAttribute("id-chat", percakapan_id);
}

let percakapan_id = "";
let dipilihJSON = {};
document.addEventListener('DOMContentLoaded', function() {
    const chatDetailMessages = document.getElementById('chat-detail-messages');
    chatDetailMessages.innerHTML = document.getElementById('loading-state-template').innerHTML;

    // Handle Conversastion
    const JsonRaw = document.getElementById('dipilih-data').textContent.trim();
    if(JsonRaw !== "null") {
        dipilihJSON = JsonRaw === "" ? {} : JSON.parse(atob(JSON.parse(JsonRaw)));        
        percakapan_id = JsonRaw === "" ? "" : dipilihJSON.id;
                    
        UpdatePercakapan(dipilihJSON);
    } else {
        chatDetailMessages.innerHTML = document.getElementById('empty-state-template').innerHTML;
    }
});

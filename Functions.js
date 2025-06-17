function formatDateWithOffset(date) {
    const pad = (n, z = 2) => ('00' + n).slice(-z);
    const offset = -date.getTimezoneOffset();
    const sign = offset >= 0 ? '+' : '-';
    const absOffset = Math.abs(offset);
    const hours = pad(Math.floor(absOffset / 60));
    const minutes = pad(absOffset % 60);
    return (
        date.getFullYear() + '-' +
        pad(date.getMonth() + 1) + '-' +
        pad(date.getDate()) + 'T' +
        pad(date.getHours()) + ':' +
        pad(date.getMinutes()) + ':' +
        pad(date.getSeconds()) + '.' +
        (date.getMilliseconds() * 1000).toString().padStart(7, '0') +
        sign + hours + ':' + minutes
    );
}

function meetNow() {
    let now = new Date();
    let endTime = new Date(now.getTime() + 30 * 60000);

    now = formatDateWithOffset(now);
    endTime = formatDateWithOffset(endTime);

    const meetingDetails = {
        subject: "Meet Now",
        start_time: now,
        end_time: endTime,
        attendees: [
            {
                emailAddress: {
                    address: "[Email Address of Attendee]",
                    name: "[Name of Attendee]"
                },
                type: "required"
            }
        ],
        organizer: "[object ID of Organizer]" // Get Object ID of user (can be grabbed if SSO is being used)
    };

    fetch('http://127.0.0.1:8080/create-event', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(meetingDetails)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    })
    .then(data => {
        console.log('Meeting created successfully:', data);
        document.getElementById('chatBox').innerHTML += `<p>${data.body}</p>`;
    })
    .catch(error => {
        console.error('Error creating meeting:', error);
    });
}

function openScheduleModal() {
    document.getElementById('modalOverlay').style.display = 'block';
    document.getElementById('scheduleModal').style.display = 'block';
}

function closeScheduleModal() {
    document.getElementById('modalOverlay').style.display = 'none';
    document.getElementById('scheduleModal').style.display = 'none';
}

function scheduleMeeting() {
    const date = document.getElementById('meetingDate').value;
    const time = document.getElementById('meetingTime').value;
    const meetingDetails = {
            subject: "Scheduled Meeting",
            start_time: new Date(`${date}T${time}`).toISOString(),
            end_time: new Date(new Date(`${date}T${time}`).getTime() + 30 * 60000).toISOString(),
            attendees: [
                {
                    emailAddress: {
                        address: "[Email Address of Attendee]",
                        name: "[Name of Attendee]"
                    },
                    type: "required"
                }
            ],
            organizer: "[object ID of Organizer]"
        };
    if (date && time) {
        fetch('http://127.0.0.1:8080/create-event', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(meetingDetails)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Meeting created successfully:', data);
            document.getElementById('chatBox').innerHTML += `<p>Meeting scheduled for ${date} at ${time}, please check your calendar to confirm reciept</p>`;
            closeScheduleModal();
        })
    } else {
        alert("Please select both date and time.");
    }
}
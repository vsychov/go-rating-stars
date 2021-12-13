import rater from 'rater-js'
import '../css/main.css';

//bg color is #f48e22
if (window.location.hash) {
    document.body.classList.add(window.location.hash.substring(1));
}

async function postData(url = '', data = {}) {
    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        redirect: 'follow',
        body: JSON.stringify(data)
    });

    return await response.json();
}

let raterAPI = rater({
    element: document.querySelector("#rater"),
    rateCallback: function rateCallback(rating, done) {
        postData('../../' + window.resourceId + '/vote/', {vote: rating})
            .then((data) => {
                if (!data.hasOwnProperty('rating') || data.rating === 0) {
                    done();
                    return;
                }
                document.querySelector("#totalVotes").innerHTML = data.totalVotes;
                this.setRating(data.rating);
                this.disable();
                done();
            });
    },
    starSize: 18,
    rating: window.rating,
    max: 5,
    step: 1,
    readOnly: window.readOnly,
});

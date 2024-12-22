fetch(`/api/users/isAdmin`)
.then(response => response.json())
.then(
    data => {
        if(data.isAdmin === true) {
            const navElement = document.querySelector('header nav');
            
            const addLink = document.createElement('a');
            addLink.href = '/add';
            addLink.innerHTML = "Add";
            navElement.appendChild(addLink);

            const userManagementLink = document.createElement('a');
            userManagementLink.href = '/user-management';
            userManagementLink.textContent = 'Users';
            navElement.appendChild(userManagementLink);
        }
    }
)

function myFunc() {
    return new Promise((resolve, reject) => {
        setTimeout(() => {
            resolve('Hello');
        }, 1000);
    }
    )
}

myFunc()
.then((data) => {
    console.log(data);
    return data;
})
.then((data) => {
    return new Promise((resolve, reject) => {
        setTimeout(() => {
            reject('World');
        }, 2000);
    })})
.catch((data) => {
    console.log(data);
});
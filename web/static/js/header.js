fetch(`/api/users/isAdmin`)
.then(response => response.json())
.then(
    data => {
        if(data.isAdmin === true) {
            const navElement = document.querySelector('header nav');
            const userManagementLink = document.createElement('a');
            userManagementLink.href = '/user-management';
            userManagementLink.textContent = 'Users';
            navElement.appendChild(userManagementLink);
        }
    }
)
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
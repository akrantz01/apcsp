import {DataService} from './DataService';

export class APIService {
    static apiURL = 'http://10.7.64.225:8080/api/';

    static async login(username, passwordHash) {
        return await fetch(this.apiURL + 'auth/login', {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            credentials: 'same-origin',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                password: passwordHash,
            }),
        })
            .then(response => response.json())
            .then(res => {
                console.log(res);
                if (res.status === 'success') {
                    DataService.saveUsername(username.toString());
                    return res.data.token.toString();
                } else {
                    console.log('not success');
                    return '';
                }
            })
            .then(token => {
                if (token !== '') {
                    console.log('token', token);
                    DataService.saveUserToken(token);
                    return true;
                } else {
                    console.log('no token');
                    return false;
                }
            })
            .catch(() => {
                console.log('could not contact server');
                return false;
            });
    }

    static async logout() {
        return await DataService.getUserToken().then(token =>
            fetch(this.apiURL + 'auth/logout', {
                method: 'GET',
                mode: 'cors',
                cache: 'no-cache',
                credentials: 'same-origin',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: token.toString(),
                },
            })
                .then(response => response.json().then(r => r.status))
                .catch(() => 'error'),
        );
    }

    static async register(name, email, username, passwordHash) {
        return await fetch(this.apiURL + 'users', {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            credentials: 'same-origin',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                password: passwordHash,
                email: email,
                name: name,
            }),
        })
            .then(response => response.json())
            .then(res => {
                console.log(res);
                return res.status === 'success';
            })
            .catch(() => {
                console.log('could not contact server');
                return false;
            });
    }

    static async deleteAccount(username) {
        return DataService.getUserToken().then(token =>
            fetch(this.apiURL + 'users/' + username, {
                method: 'DELETE',
                mode: 'cors',
                cache: 'no-cache',
                credentials: 'same-origin',
                headers: {
                    Accept: 'application/json',
                    Authorization: token,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    user: username,
                }),
            }).then(response => {
                console.log(response);
                return response;
            }),
        );
    }

    static async getUserChats() {
        return DataService.getUserToken().then(token =>
            fetch(this.apiURL + 'api/chats', {
                method: 'GET',
                mode: 'cors',
                cache: 'no-cache',
                credentials: 'same-origin',
                headers: {
                    Accept: 'application/json',
                    Authorization: token,
                    'Content-Type': 'application/json',
                },
            })
                .then(res => res.json())
                .then(response => {
                    if (response.status === 'success') {
                        return response.data;
                    } else {
                        return 'error';
                    }
                }),
        );
    }
}

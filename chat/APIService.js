import {DataService} from './DataService';

export class APIService {
    static apiURL = 'http://192.168.0.111:8080/api/';

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
                if (res.status === 'success') {
                    if (res.hasOwnProperty('data')) {
                        return res.data.token.toString();
                    } else {
                        console.log('success but corrupted payload');
                        return '';
                    }
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
                console.log('catch');
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
            }).then(response => response.json().then(r => r.status)),
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
}

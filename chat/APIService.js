import AsyncStorage from '@react-native-community/async-storage';

export class APIService {
    static apiURL = 'https://akrantz01.github.io/apcsp/api/';

    static async login(username, passwordHash) {
        await fetch(this.apiURL + 'auth/login', {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                password: passwordHash,
            }),
        }).then(response => {
            AsyncStorage.setItem('authToken', response.data.token);
        });
    }

    static async logout() {
        await fetch(this.apiURL + 'auth/logout', {
            method: 'GET',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
                Authorization: AsyncStorage.getItem('authKey'),
            },
        }).then(response => {
            return response;
        });
    }

    static async register(name, email, username, passwordHash) {
        await fetch(this.apiURL + 'users', {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                password: passwordHash,
                email: email,
                name: name,
            }),
        }).then(response => {
            return response;
        });
    }
}

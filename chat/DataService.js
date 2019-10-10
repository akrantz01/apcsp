import AsyncStorage from '@react-native-community/async-storage';

export class DataService {
    static async isLoggedIn() {
        const userToken = await AsyncStorage.getItem('authToken');
        return !(!userToken || userToken === '');
    }

    static async saveUserToken(token) {
        AsyncStorage.setItem('authToken', token);
    }

    static async getUserToken() {
        return await AsyncStorage.getItem('authToken', '');
    }

    static async removeUserToken() {
        AsyncStorage.setItem('authToken', '');
    }

    static async saveUsername(username) {
        AsyncStorage.setItem('currentUser', username);
    }

    static async getUsername() {
        return await AsyncStorage.getItem('currentUser', '');
    }

    static getSavedMessages(id) {
        console.log(id);
        let messages = AsyncStorage.getItem('messages_' + id);
        return messages === null
            ? messages
            : [
                  {
                      _id: 1,
                      text: 'Hello developer',
                      createdAt: new Date(),
                      user: {
                          _id: 2,
                          name: 'React Native',
                          avatar: 'https://placeimg.com/140/140/any',
                      },
                  },
              ];
    }

    static saveMessage(id, messages) {
        AsyncStorage.setItem('messages_' + id, messages);
    }

    static getAvatar() {}
}

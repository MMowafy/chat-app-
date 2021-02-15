class App < ApplicationRecord
    has_many :chats
    before_create :generate_token

      protected

      def generate_token
        self.token = loop do
          random_token = SecureRandom.urlsafe_base64
          break random_token unless App.exists?(token: random_token)
        end
      end
end

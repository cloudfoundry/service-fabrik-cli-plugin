package helper

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/SAP/service-fabrik-cli-plugin/errors"
	"io/ioutil"
	"testing"
)

func TestConfigJsonParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Json Parser Test Suite")
}

var _ = Describe("helper", func() {
	Context("Finding token info", func() {
		It("Token info should match", func() {
			accessToken := "bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImxlZ2FjeS10b2tlbi1rZXkiLCJ0eXAiOiJKV1QifQ.eyJqdGkiOiI0ZTA3MDMxYTJhYTA0ZTc4YTc4NTc3N2IxY2MyY2NjYyIsInN1YiI6IjNlZjIzNmQ2LWZmZWMtNGY3Yi1iMzQxLWIwZWI1MmMyMDZiNyIsInNjb3BlIjpbImNsb3VkX2NvbnRyb2xsZXIuYWRtaW4iLCJjbG91ZF9jb250cm9sbGVyLndyaXRlIiwiZG9wcGxlci5maXJlaG9zZSIsIm9wZW5pZCIsInNjaW0ucmVhZCIsInVhYS51c2VyIiwiY2xvdWRfY29udHJvbGxlci5yZWFkIiwicGFzc3dvcmQud3JpdGUiLCJzY2ltLndyaXRlIl0sImNsaWVudF9pZCI6ImNmIiwiY2lkIjoiY2YiLCJhenAiOiJjZiIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiIzZWYyMzZkNi1mZmVjLTRmN2ItYjM0MS1iMGViNTJjMjA2YjciLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX25hbWUiOiJhZG1pbiIsImVtYWlsIjoiYWRtaW4iLCJhdXRoX3RpbWUiOjE1MDA1NDczNzMsInJldl9zaWciOiJjNjgxODE5OCIsImlhdCI6MTUwMDU0NzM3MywiZXhwIjoxNTAwNTQ3OTczLCJpc3MiOiJodHRwczovL3VhYS5jZi5zZi1kZXY0LmF3cy5zYXBjbG91ZC5pby9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjbG91ZF9jb250cm9sbGVyIiwic2NpbSIsInBhc3N3b3JkIiwiY2YiLCJ1YWEiLCJvcGVuaWQiLCJkb3BwbGVyIl19.RzPnFW8pfMhn0V0fm0RKT8-fXPBhLshHOBew3_jd3TEx9sOlrpeCKIARO1hYLiCORKM3pLulYGPQ-9SHEgYwmNNARvbOhSsMNXbStZZ2nElGuiOIJgrK5ItwkVGeC_etdt_tXJ3MNj-wzFhGHvPCVbMn7R-gXu_oPr3nJC2PGc5H4ArO32hTcVEcCCKohWeJFvCPqYaIzTHll3q2tSX3WXwfr3bg5WHfBrcKkfMT_WX2D5wuOJzoD4WbTG_fBeNUFrBG_PsWWMX765zF6b-bC56XnU3WmxlOyQAhBpPH2UYRTfKEhnN2E2KYP2Bt753wwTFsEnZsG7g1XrfrlVX03w"
			tokenInfo := NewTokenInfo(accessToken)
			Expect(tokenInfo.UserGUID).To(Equal("3ef236d6-ffec-4f7b-b341-b0eb52c206b7"))
		})
	})
	Context("Finding decode access token", func() {
		It("Decoded access token should match", func() {
			accessToken := "bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImxlZ2FjeS10b2tlbi1rZXkiLCJ0eXAiOiJKV1QifQ.eyJqdGkiOiI0ZTA3MDMxYTJhYTA0ZTc4YTc4NTc3N2IxY2MyY2NjYyIsInN1YiI6IjNlZjIzNmQ2LWZmZWMtNGY3Yi1iMzQxLWIwZWI1MmMyMDZiNyIsInNjb3BlIjpbImNsb3VkX2NvbnRyb2xsZXIuYWRtaW4iLCJjbG91ZF9jb250cm9sbGVyLndyaXRlIiwiZG9wcGxlci5maXJlaG9zZSIsIm9wZW5pZCIsInNjaW0ucmVhZCIsInVhYS51c2VyIiwiY2xvdWRfY29udHJvbGxlci5yZWFkIiwicGFzc3dvcmQud3JpdGUiLCJzY2ltLndyaXRlIl0sImNsaWVudF9pZCI6ImNmIiwiY2lkIjoiY2YiLCJhenAiOiJjZiIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiIzZWYyMzZkNi1mZmVjLTRmN2ItYjM0MS1iMGViNTJjMjA2YjciLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX25hbWUiOiJhZG1pbiIsImVtYWlsIjoiYWRtaW4iLCJhdXRoX3RpbWUiOjE1MDA1NDczNzMsInJldl9zaWciOiJjNjgxODE5OCIsImlhdCI6MTUwMDU0NzM3MywiZXhwIjoxNTAwNTQ3OTczLCJpc3MiOiJodHRwczovL3VhYS5jZi5zZi1kZXY0LmF3cy5zYXBjbG91ZC5pby9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjbG91ZF9jb250cm9sbGVyIiwic2NpbSIsInBhc3N3b3JkIiwiY2YiLCJ1YWEiLCJvcGVuaWQiLCJkb3BwbGVyIl19.RzPnFW8pfMhn0V0fm0RKT8-fXPBhLshHOBew3_jd3TEx9sOlrpeCKIARO1hYLiCORKM3pLulYGPQ-9SHEgYwmNNARvbOhSsMNXbStZZ2nElGuiOIJgrK5ItwkVGeC_etdt_tXJ3MNj-wzFhGHvPCVbMn7R-gXu_oPr3nJC2PGc5H4ArO32hTcVEcCCKohWeJFvCPqYaIzTHll3q2tSX3WXwfr3bg5WHfBrcKkfMT_WX2D5wuOJzoD4WbTG_fBeNUFrBG_PsWWMX765zF6b-bC56XnU3WmxlOyQAhBpPH2UYRTfKEhnN2E2KYP2Bt753wwTFsEnZsG7g1XrfrlVX03w"
			decodedToken, _ := DecodeAccessToken(accessToken)
			Expect(string(decodedToken)).To(Equal(`{"jti":"4e07031a2aa04e78a785777b1cc2cccc","sub":"3ef236d6-ffec-4f7b-b341-b0eb52c206b7","scope":["cloud_controller.admin","cloud_controller.write","doppler.firehose","openid","scim.read","uaa.user","cloud_controller.read","password.write","scim.write"],"client_id":"cf","cid":"cf","azp":"cf","grant_type":"password","user_id":"3ef236d6-ffec-4f7b-b341-b0eb52c206b7","origin":"uaa","user_name":"admin","email":"admin","auth_time":1500547373,"rev_sig":"c6818198","iat":1500547373,"exp":1500547973,"iss":"https://uaa.cf.sf-dev4.aws.sapcloud.io/oauth/token","zid":"uaa","aud":["cloud_controller","scim","password","cf","uaa","openid","doppler"]}`))
		})
	})
	Context("Getting access token", func() {
		It("Access token should match", func() {
			file, err := ioutil.ReadFile("../test/config_test.json")
			if err != nil {
				errors.FileReadingError("config.json")
			}
			accessToken := GetAccessToken(file)
			Expect(accessToken).To(Equal("bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImxlZ2FjeS10b2tlbi1rZXkiLCJ0eXAiOiJKV1QifQ.eyJqdGkiOiJlYTI3YzZhNmRiZDg0Y2Y0YmYxNWExZTM2MzFmOGQwOSIsInN1YiI6IjVmNTBlZjdlLWU3NDUtNDIwZC04NTQ2LWM5OTEwZWZhOWUxYyIsInNjb3BlIjpbImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSIsIm9wZW5pZCIsImRvcHBsZXIuZmlyZWhvc2UiLCJzY2ltLndyaXRlIiwic2NpbS5yZWFkIiwiY2xvdWRfY29udHJvbGxlci5hZG1pbiIsInVhYS51c2VyIl0sImNsaWVudF9pZCI6ImNmIiwiY2lkIjoiY2YiLCJhenAiOiJjZiIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiI1ZjUwZWY3ZS1lNzQ1LTQyMGQtODU0Ni1jOTkxMGVmYTllMWMiLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX25hbWUiOiJhZG1pbiIsImVtYWlsIjoiYWRtaW4iLCJyZXZfc2lnIjoiYTMwZWI5ZjEiLCJpYXQiOjE0ODA0MDY1MTUsImV4cCI6MTQ4MDQwNzExNSwiaXNzIjoiaHR0cHM6Ly91YWEuY2Yuc2VydmljZS1mYWJyaWsuc2M2LnNhcGNsb3VkLmlvL29hdXRoL3Rva2VuIiwiemlkIjoidWFhIiwiYXVkIjpbInNjaW0iLCJjbG91ZF9jb250cm9sbGVyIiwicGFzc3dvcmQiLCJjZiIsInVhYSIsIm9wZW5pZCIsImRvcHBsZXIiXX0.q_AYTPgwR5VP6-i3QNbTaATKCsDKe8B76udSkqCKxa-ZLjUFppogONP2Yd6S_f_mG23SjyUYVnYay2d62W1I4i9Ih28aBcYqaCsVyecvenr3ujS_P4KTjTfrm7wq-qdIF4H0DGNCreNU_XImzDug8dGvJFen9duGHsgNjcUGY7g"))
		})
	})
	Context("Getting space guid", func() {
		It("Space guid should match", func() {
			file, err := ioutil.ReadFile("../test/config_test.json")
			if err != nil {
				errors.FileReadingError("config.json")
			}
			spaceGuid := GetSpaceGUID(file)
			Expect(spaceGuid).To(Equal("b0728cce-2eef-4a8b-ac57-b480f2c48461"))
		})
	})
	Context("Getting space name", func() {
		It("Space name should match", func() {
			file, err := ioutil.ReadFile("../test/config_test.json")
			if err != nil {
				errors.FileReadingError("config.json")
			}
			spaceName := GetSpaceName(file)
			Expect(spaceName).To(Equal("postgresql_test"))
		})
	})
	Context("Getting org name", func() {
		It("Org name should match", func() {
			file, err := ioutil.ReadFile("../test/config_test.json")
			if err != nil {
				errors.FileReadingError("config.json")
			}
			orgName := GetOrgName(file)
			Expect(orgName).To(Equal("test"))
		})
	})
	Context("Getting api endpoint", func() {
		It("API endpoint should match", func() {
			file, err := ioutil.ReadFile("../test/config_test.json")
			if err != nil {
				errors.FileReadingError("config.json")
			}
			apiEndpoint := GetApiEndpoint(file)
			Expect(apiEndpoint).To(Equal("https://api.cf.service-fabrik.io"))
		})
	})
})

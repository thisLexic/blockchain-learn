package main

import (
	"strconv"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	peer "github.com/hyperledger/fabric-protos-go/peer"
)

type TokenChaincode struct {
}

type Bid struct {
	Company string
	Value uint64
}

type Project struct {
	Name string
	Bids []Bid
}

type ProjectList struct {
	ProjectList []Project
}

func (token *TokenChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	var bid1 Bid = Bid{Company:"Com 1", Value: 110}
	var bid2 Bid = Bid{Company:"Com 2", Value: 220}
	var bidSlice []Bid = []Bid{bid1, bid2}
	var project1 Project = Project{Name:"Proj A",Bids:bidSlice}

	var bid3 Bid = Bid{Company:"Com 3", Value: 330}
	var bid4 Bid = Bid{Company:"Com 4", Value: 440}
	var bid5 Bid = Bid{Company:"Com 5", Value: 550}
	bidSlice = []Bid{bid3, bid4, bid5}
	var project2 Project = Project{Name:"Proj B",Bids:bidSlice}

	var projects []Project = []Project{project1,project2}
	var projectList ProjectList = ProjectList{ProjectList: projects}
	jsonProjectList, _ := json.Marshal(projectList)


	stub.PutState("governmentAgencies", []byte(`{"DPWH"}`))
	stub.PutState("DPWH", jsonProjectList)
	return shim.Success([]byte("true"))
}

func (token *TokenChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	funcName, args := stub.GetFunctionAndParameters()

	if funcName == "getProjects" {
		return getProjects(stub, args)
	} else if funcName == "addProject" {
		return addProject(stub, args)
	} else if funcName == "addBid" {
		return addBid(stub, args)
	}

	return shim.Error(funcName + ` is an invalid function!`)
}


func getProjects(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	governmentAgency := args[0]

	projects, _ := stub.GetState(governmentAgency)
	return shim.Success(projects)
}


func addProject(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	governmentAgency := args[0]
	projectName := args[1]

	projectListBytes, _ := stub.GetState(governmentAgency)

	var projectList ProjectList
	_ = json.Unmarshal(projectListBytes, &projectList)

	var project Project = Project{Name: projectName}
	projectList.ProjectList = append(projectList.ProjectList, project)

	jsonProjectList, _ := json.Marshal(projectList)
	stub.PutState(governmentAgency, jsonProjectList)

	return shim.Success(jsonProjectList)	
}


// Fix this part
func addBid(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	governmentAgency := args[0]
	projectName := args[1]
	bidCompany := args[2]
	bid := args[3]

	projectListBytes, _ := stub.GetState(governmentAgency)

	var projectList ProjectList
	_ = json.Unmarshal(projectListBytes, &projectList)

	var queriedProject *Project
	for i, project := range projectList.ProjectList {
		if project.Name == projectName {
			queriedProject = &projectList.ProjectList[i]
		}
	}
	if queriedProject == nil {
		return shim.Error(projectName + " is not a project of " + governmentAgency)
	}

	bidInt, _ := strconv.ParseInt(bid, 10, 64)
	bidUint := uint64(bidInt)
	var bidObj Bid = Bid{Company: bidCompany, Value: bidUint}
	queriedProject.Bids = append(queriedProject.Bids, bidObj)

	jsonProject, _ := json.Marshal(queriedProject)
	jsonProjectList, _ := json.Marshal(projectList)
	stub.PutState(governmentAgency, jsonProjectList)

	return shim.Success(jsonProject)
}


func main() {
	err := shim.Start(new(TokenChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}